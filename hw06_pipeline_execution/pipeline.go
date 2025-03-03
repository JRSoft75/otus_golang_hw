package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	// Создаем пайплайн, последовательно соединяя все стейджи.
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		out = stage(doFunction(out, done))
	}
	return doFunction(out, done)
}

func doFunction(inputCh In, stop In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out) // Закрываем выходной канал после завершения работы стейджа.
			for skip := range inputCh {
				_ = skip // очищаем канал от данных
			}
		}()
		for {
			select {
			case <-stop: // Проверяем сигнал остановки.
				return
			case value, ok := <-inputCh:
				if !ok {
					return // Входной канал закрыт
				}
				out <- value
			}
		}
	}()
	return out
}
