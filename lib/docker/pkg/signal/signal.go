package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var SignalMap = map[string]syscall.Signal{
	"ABRT": syscall.SIGABRT,
	"ALRM": syscall.SIGALRM,
	"BUG":  syscall.SIGBUS,
	"CHLD": syscall.SIGCHLD,
	"CONT": syscall.SIGCONT,
	//"EMT":    syscall.SIGEMT,
	"FPE": syscall.SIGFPE,
	"HUP": syscall.SIGHUP,
	"ILL": syscall.SIGILL,
	//"INFO":   syscall.SIGINFO,
	"INT":    syscall.SIGINT,
	"IO":     syscall.SIGIO,
	"IOT":    syscall.SIGIOT,
	"KILL":   syscall.SIGKILL,
	"PIPE":   syscall.SIGPIPE,
	"PROF":   syscall.SIGPROF,
	"QUIT":   syscall.SIGQUIT,
	"SEGV":   syscall.SIGSEGV,
	"STOP":   syscall.SIGSTOP,
	"SYS":    syscall.SIGSYS,
	"TERM":   syscall.SIGTERM,
	"TRAP":   syscall.SIGTRAP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"URG":    syscall.SIGURG,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"VTALRM": syscall.SIGVTALRM,
	"WINCH":  syscall.SIGWINCH,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
}

func CatchAll(sigc chan os.Signal) {
	handledSigs := []os.Signal{}
	for _, s := range SignalMap {
		handledSigs = append(handledSigs, s)
	}
	signal.Notify(sigc, handledSigs...)
}

func StopCatch(sigc chan os.Signal) {
	signal.Stop(sigc)
	close(sigc)
}
