package service

// PomodoroConfig define el patrón clásico (o de pruebas).
type PomodoroConfig struct {
	PomodorosPerCycle int
	ShortBreakMinutes int
	LongBreakMinutes  int
	FocusMinutes      int
}

// MODO PRUEBAS
// - 4 pomodoros por ciclo (igual que el real)
// - Cada pomodoro de 1 minuto
// - Descansos cortos de 1 minuto
// - Descanso largo (después del 4.º) de 2 minutos
var DefaultPomodoroConfig = PomodoroConfig{
	PomodorosPerCycle: 4, // 4 focos por ciclo
	ShortBreakMinutes: 1, // 1 min de descanso corto
	LongBreakMinutes:  2, // 2 min de descanso largo
	FocusMinutes:      1, // 1 min de foco
}

// PomodoroCycleState representa en qué parte del ciclo estás.
type PomodoroCycleState struct {
	TotalPomodoros int  // pomodoros totales (en la vida de la tarea)
	IndexInCycle   int  // 1..4 dentro del ciclo actual
	IsCycleEnd     bool // true si se acaba de terminar el 4to
	NextBreakMin   int  // 1 o 2 (según corto/largo en modo test)
	CyclesDone     int  // cuántos ciclos completos (de 4) llevas
}

// ComputePomodoroState calcula el estado después de terminar un foco.
func ComputePomodoroState(totalPomodorosAfter int, cfg PomodoroConfig) PomodoroCycleState {
	index := totalPomodorosAfter % cfg.PomodorosPerCycle
	if index == 0 {
		index = cfg.PomodorosPerCycle
	}

	isEnd := index == cfg.PomodorosPerCycle

	nextBreak := cfg.ShortBreakMinutes
	if isEnd {
		nextBreak = cfg.LongBreakMinutes
	}

	cyclesDone := totalPomodorosAfter / cfg.PomodorosPerCycle

	return PomodoroCycleState{
		TotalPomodoros: totalPomodorosAfter,
		IndexInCycle:   index,
		IsCycleEnd:     isEnd,
		NextBreakMin:   nextBreak,
		CyclesDone:     cyclesDone,
	}
}
