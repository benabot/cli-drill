"""Unix PTY regression checks. Pass the path to a built cli-drill binary."""

import contextlib
import fcntl
import json
import os
from pathlib import Path
import pty
import select
import signal
import struct
import subprocess
import sys
import tempfile
import termios
import time


class Capture:
    def __init__(self, process, master, slave, progress, chapter_id):
        self.process = process
        self.master = master
        self.slave = slave
        self.progress = progress
        self.chapter_id = chapter_id
        self.pending = b""
        self.probes = [
            (b"\x1b]11;?\x1b\\", b"\x1b]11;rgb:0000/0000/0000\x1b\\"),
            (b"\x1b[6n", b"\x1b[1;1R"),
        ]

    def pump(self):
        if select.select([self.master], [], [], 0.02)[0]:
            self.pending += os.read(self.master, 8192)
        for request, response in self.probes:
            if request in self.pending:
                self.pending = self.pending.replace(request, b"", 1)
                os.write(self.master, response)

    def wait(self, condition):
        deadline = time.monotonic() + 5
        while not condition():
            if time.monotonic() >= deadline:
                raise AssertionError(f"PTY timeout: {self.pending[-1500:]!r}")
            self.pump()

    def expect(self, value):
        value = value.encode()
        self.wait(lambda: value in self.pending)
        self.pending = self.pending.split(value, 1)[1]

    def is_raw(self):
        return not bool(termios.tcgetattr(self.slave)[3] & termios.ICANON)

    def send(self, value, raw=True):
        self.wait(lambda: self.is_raw() == raw)
        os.write(self.master, value)

    def finish(self):
        self.wait(lambda: self.process.poll() is not None)
        self.pump()
        return self.process.returncode

    def assert_score(self, attempts, correct):
        if attempts == 0:
            assert not self.progress.exists(), self.progress
            return
        state = json.loads(self.progress.read_text())
        assert state["version"] == 1, state
        assert state["chapter_scores"][self.chapter_id] == {
            "attempts": attempts, "correct": correct,
        }, state


@contextlib.contextmanager
def session(binary, difficulty="beginner", all_exercises=False, chapter_id="01-raccourcis-terminal"):
    with tempfile.TemporaryDirectory(prefix="cli-drill-capture-") as temp:
        progress = Path(temp) / "progress.json"
        master, slave = pty.openpty()
        fcntl.ioctl(slave, termios.TIOCSWINSZ, struct.pack("HHHH", 24, 80, 0, 0))
        before = termios.tcgetattr(slave)
        args = [binary, "--progress", str(progress)]
        if difficulty is not None:
            args.extend(["train", chapter_id, "--difficulty", difficulty])
        if all_exercises:
            args.append("--all")
        process = subprocess.Popen(
            args,
            stdin=slave, stdout=slave, stderr=slave,
            env={"XDG_DATA_HOME": temp, "TERM": "xterm-256color", "NO_COLOR": "1"},
        )
        capture = Capture(process, master, slave, progress, chapter_id)
        try:
            yield capture
            after = termios.tcgetattr(slave)
            # Darwin may mark pending input for retyping after a mode change.
            # PENDIN is a transient state flag, not a terminal setting.
            for attrs in (before, after):
                attrs[3] &= ~getattr(termios, "PENDIN", 0)
            assert before == after, (before, after, capture.pending)
        finally:
            if process.poll() is None:
                process.kill()
                process.wait()
            termios.tcsetattr(slave, termios.TCSANOW, before)
            os.close(master)
            os.close(slave)


def check_signals(binary):
    for sig in (signal.SIGINT, signal.SIGTERM, signal.SIGHUP, signal.SIGQUIT, signal.SIGTSTP):
        with session(binary) as capture:
            capture.expect("Waiting for key")
            capture.send(b"\x02")
            capture.expect("Correct.")
            capture.wait(capture.is_raw)
            capture.process.send_signal(sig)
            capture.expect("capture clavier interrompue")
            assert capture.finish() == 128 + sig
            capture.assert_score(1, 1)
        print(f"OK {sig.name}: terminal restored and saved answer preserved")


def check_escape_and_labels(binary):
    for exit_key in (b"\x1b", b"\x03"):
        with session(binary) as capture:
            capture.expect("Waiting for key")
            capture.send(b"\x1b[A\x1bOP\x1b[200~test\x1b[201~")
            capture.expect("Waiting for key")
            capture.assert_score(0, 0)
            for key, label in ((b"\x02", "Ctrl+B"), (b"\x06", "Ctrl+F"),
                               (b"\x7f", "DEL (Retour arrière)"), (b"\t", "Tab (Ctrl+I)")):
                capture.send(key)
                capture.expect("Recu: " + label)
                capture.expect("Correct.")
                capture.send(b"\r")
                capture.expect("Waiting for key")
            capture.send(exit_key)
            capture.expect("Session interrompue")
            assert capture.finish() == 0
            capture.assert_score(4, 4)
        print(f"OK {exit_key!r}: unsupported sequences and exit do not count")


def check_horizontal_arrows(binary):
    for prefix in (b"\x1b[", b"\x1bO"):
        with session(binary) as capture:
            for key, label in ((b"D", "←"), (b"C", "→")):
                capture.expect("Waiting for key")
                capture.send(prefix + key)
                capture.expect("Recu: " + label)
                capture.expect("Correct.")
                capture.send(b"\r")
            capture.expect("Waiting for key")
            capture.send(prefix + b"D")
            capture.expect("Waiting for key")
            capture.send(b"\x1b")
            capture.expect("Session interrompue")
            capture.expect("Bonnes réponses : 2/2")
            assert capture.finish() == 0
            capture.assert_score(2, 2)
    print("OK capture: CSI and SS3 arrows accepted on matching exercises")


def check_mixed_session(binary):
    keys = (b"\x04", b"\x14", b"\x10", b"\x0e", b"\x1f", b"\x07")
    with session(binary, "intermediate", all_exercises=True) as capture:
        for key in (b"x", *keys[1:]):
            capture.expect("Waiting for key")
            capture.send(key)
            if key == b"\x1f":
                capture.expect("Recu: Ctrl+_")
            capture.expect("Pas encore." if key == b"x" else "Correct.")
            capture.send(b"\r")
        capture.expect("Exercice 7/11")
        capture.send(b"\x1b\r", raw=False)
        capture.expect("Session interrompue")
        capture.expect("Bonnes réponses : 5/6")
        assert capture.finish() == 0
        assert b"r Revoir mes erreurs" not in capture.pending
        capture.assert_score(6, 5)
    print("OK mixed session: text escape exits without opening the review")

    with session(binary, "intermediate", all_exercises=True) as capture:
        for key in keys:
            capture.expect("Waiting for key")
            capture.send(key)
            capture.expect("Correct.")
            capture.send(b"\r")
        for index, answer in enumerate((
            "git switch |sauvegarde", "git status| --short", "grep | serveur.log",
            "Interrompre la saisie avec Ctrl+C.",
            "Ajouter n au préfixe puis redemander la complétion.",
        ), 7):
            capture.expect(f"Exercice {index}/11")
            capture.send((answer + "\r").encode(), raw=False)
            capture.expect("Correct.")
        capture.expect("Bonnes réponses : 11/11")
        assert capture.finish() == 0
        assert b"Bilan du lot" not in capture.pending
        capture.assert_score(11, 11)
    print("OK mixed session: one score covers both formats")


def check_short_mixed_visit(binary):
    with session(binary, "intermediate") as capture:
        for index, key in enumerate((b"x", b"\x14", b"\x10", b"\x0e", b"\x1f", b"\x07"), 1):
            capture.expect(f"Progression: {index}/6")
            capture.expect("Waiting for key")
            capture.send(key)
            capture.expect("Pas encore." if key == b"x" else "Correct.")
            capture.send(b"\r")
        capture.expect("Bilan du lot")
        capture.expect("Bonnes réponses : 5/6")
        capture.expect("r Revoir mes erreurs")
        capture.send(b"r")
        capture.expect("Progression: 1/1")
        capture.expect("Waiting for key")
        capture.send(b"\x04")
        capture.expect("Correct.")
        capture.send(b"\r")
        capture.expect("Bilan de la reprise")
        capture.expect("Bonnes réponses : 1/1")
        capture.expect("Lot initial : 5/6")
        capture.expect("c Continuer")
        capture.send(b"c")
        capture.expect("Exercice 1/5")
        capture.send(b"git switch |sauvegarde\r", raw=False)
        capture.expect("Correct.")
        capture.expect("Exercice 2/5")
        capture.send(b"\x1b\r", raw=False)
        capture.expect("Session interrompue")
        assert capture.finish() == 0
        capture.assert_score(8, 7)
    print("OK short visit: six captures, one retry, then the text batch")


def check_tui_return(binary):
    with session(binary, difficulty=None) as capture:
        capture.expect("Start training")
        capture.send(b"\x1b[B\x1b[B\r")
        capture.expect("Waiting for key")
        capture.send(b"\x02")
        capture.expect("Correct.")
        capture.send(b"\x1b")
        capture.expect("Start training")
        capture.send(b"q")
        assert capture.finish() == 0
        capture.assert_score(1, 1)
    print("OK TUI: escape returns to the home screen and preserves the answer")


def check_tui_full_count(binary):
    with session(binary, difficulty=None, chapter_id="00-premiers-pas") as terminal:
        terminal.expect("Start training")
        terminal.send(b"\x1b[B\r")
        terminal.expect("/ search")
        terminal.send(b"/")
        terminal.expect("type to filter")
        terminal.send(b"00-premiers-pas")
        terminal.expect("00-premiers-pas")
        terminal.send(b"\r")
        terminal.expect("esc clear filter")
        terminal.send(b"\r")
        terminal.expect("Ma séance")
        terminal.send(b"a")
        terminal.expect("Parcours intégral (12)")
        terminal.send(b"\r")
        terminal.expect("Exercice 1/12")
        terminal.send(b"pwd\r")
        terminal.expect("Correct.")
        terminal.send(b"\x03")
        assert terminal.finish() == 0
        terminal.assert_score(1, 1)
    print("OK TUI: full beginner session displays Exercice 1/12")


def check_text_escape(binary):
    with session(binary, chapter_id="00-premiers-pas") as terminal:
        terminal.expect("Exercice 1/6")
        terminal.send(b"pwd\r", raw=False)
        terminal.expect("Correct.")
        terminal.expect("Exercice 2/6")
        for line in (b"\x1b[D", b"\x1b[C", b"\x1b[A", b"\x1b[B", b"\x1bOP", b"l\x1b[Ds", b""):
            terminal.send(line + b"\r", raw=False)
            terminal.expect("Exercice 2/6")
            terminal.assert_score(1, 1)
        terminal.send(b"ls\r", raw=False)
        terminal.expect("Correct.")
        terminal.expect("Exercice 3/6")
        terminal.send(b"\x1b\r", raw=False)
        terminal.expect("Session interrompue.")
        terminal.expect("Bonnes réponses : 2/2")
        assert terminal.finish() == 0
        terminal.assert_score(2, 2)
    print("OK text CLI: arrow sequences ignored, standalone Escape exits with saved answers")


if __name__ == "__main__":
    binary = str(Path(sys.argv[1]).resolve())
    check_text_escape(binary)
    check_signals(binary)
    check_escape_and_labels(binary)
    check_horizontal_arrows(binary)
    check_mixed_session(binary)
    check_short_mixed_visit(binary)
    check_tui_return(binary)
    check_tui_full_count(binary)
