package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	//wg := sync.WaitGroup{}
	//currentOut := in // Начинаем с входного канала.
	// Создаем выходной канал
	out := make(Bi)

	// Запускаем горутину для обработки стейджей
	go func() {
		defer close(out) // Закрываем выходной канал после завершения

		// Создаем промежуточные каналы для каждого стейджа
		channels := make([]Bi, len(stages))
		for i := range channels {
			channels[i] = make(Bi)
		}

		// Запускаем каждый стейдж в отдельной горутине
		for i, stage := range stages {
			//wg.Add(1)
			go func(i int, stage Stage, in2 In) {
				defer func() {
					close(channels[i]) // Закрываем канал после завершения стейджа
					//wg.Done()
				}()
				for {
					select {
					case <-done:
						return // Завершаем, если канал done закрыт
					default:

						//case v, ok := <-in2:
						//	if !ok {
						//		return // Входной канал закрыт
						//	}
						channels[i] <- stage(in2)
						//result := stage(makeOneElemChan(v)) // Передаем данные через одноразовый канал.
						//for r := range result {
						//	select {
						//	case channels[i] <- r: // Отправляем результат в выходной канал.
						//	}
						//}

					}
				}
			}(i, stage, in)
			in = channels[i] // Передаем выходной канал текущего стейджа как входной для следующего
		}

		// Передаем данные из последнего стейджа в выходной канал
		go func() {
			for v := range in {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}()
	}()

	//// Создаем пайплайн, последовательно соединяя все стейджи.
	//for _, stage := range stages {
	//	inputNext := currentOut // Сохраняем текущий выход как вход для следующего стейджа.
	//
	//	// Создаем временный канал
	//	tmpOut := make(chan interface{})
	//
	//	go func(s Stage, input In, output chan interface{}, stop In) {
	//		defer close(output) // Закрываем выходной канал после завершения работы стейджа.
	//		for {
	//			select {
	//			case <-stop: // Проверяем сигнал остановки.
	//				return
	//			case inputValue, ok := <-input:
	//				if !ok {
	//					return // Входной канал закрыт
	//				}
	//				result := s(makeOneElemChan(inputValue)) // Передаем данные через одноразовый канал.
	//				for r := range result {
	//					select {
	//					case output <- r: // Отправляем результат в выходной канал.
	//					}
	//				}
	//			}
	//		}
	//		//for inputValue := range input {
	//		//	select {
	//		//	case <-stop: // Проверяем сигнал остановки.
	//		//		return
	//		//	default:
	//		//		// Применяем стейдж к данным.
	//		//		result := s(makeOneElemChan(inputValue)) // Передаем данные через одноразовый канал.
	//		//		for r := range result {
	//		//			select {
	//		//			case <-stop: // Проверяем сигнал остановки перед отправкой данных.
	//		//				return
	//		//			case output <- r: // Отправляем результат в выходной канал.
	//		//			}
	//		//		}
	//		//	}
	//		//}
	//
	//		//for {
	//		//	select {
	//		//	case <-stop:
	//		//		return // Завершаем, если канал done закрыт
	//		//	case inputValue, ok := <-input:
	//		//		if !ok {
	//		//			return // Входной канал закрыт
	//		//		}
	//		//		select {
	//		//		case <-stop:
	//		//			return // Завершаем, если канал done закрыт
	//		//		default:
	//		//			stageResult := s(makeOneElemChan(inputValue)) // Передаем данные через одноразовый канал.
	//		//			for r := range stageResult {
	//		//				select {
	//		//				case <-stop: // Проверяем сигнал остановки перед отправкой данных.
	//		//					return
	//		//				case output <- r: // Отправляем результат в выходной канал.
	//		//				}
	//		//			}
	//		//
	//		//		}
	//		//	}
	//		//}
	//
	//		//for inputValue := range input {
	//		//	select {
	//		//	case <-stop: // Проверяем сигнал остановки.
	//		//		return
	//		//	default:
	//		//		// Применяем стейдж к данным.
	//		//		result := s(makeOneElemChan(inputValue)) // Передаем данные через одноразовый канал.
	//		//		for r := range result {
	//		//			select {
	//		//			case <-stop: // Проверяем сигнал остановки перед отправкой данных.
	//		//				return
	//		//			case output <- r: // Отправляем результат в выходной канал.
	//		//			}
	//		//		}
	//		//	}
	//		//}
	//	}(stage, inputNext, tmpOut, done)
	//
	//	currentOut = tmpOut
	//}
	//wg.Wait()
	return out
}

// Вспомогательная функция для создания канала с одним элементом.
func makeOneElemChan(v interface{}) In {
	ch := make(chan interface{})
	go func() {
		//defer close(ch)
		select {
		case ch <- v: // Отправляем значение в канал.
		}
	}()
	return ch
}

//func ExecutePipeline(in In, done In, stages ...Stage) Out {
//	// Создаем выходной канал
//	out := make(Bi)
//
//	// Запускаем горутину для обработки стейджей
//	go func() {
//		defer close(out) // Закрываем выходной канал после завершения
//
//		// Создаем промежуточные каналы для каждого стейджа
//		channels := make([]Bi, len(stages))
//		for i := range channels {
//			channels[i] = make(Bi)
//		}
//
//		// Запускаем каждый стейдж в отдельной горутине
//		for i, stage := range stages {
//			go func(i int, stage Stage, in In) {
//				defer close(channels[i]) // Закрываем канал после завершения стейджа
//				for {
//					select {
//					case v, ok := <-in:
//						if !ok {
//							return // Входной канал закрыт
//						}
//						select {
//						case <-done:
//							return // Завершаем, если канал done закрыт
//						case channels[i] <- v:
//						}
//					case <-done:
//						return // Завершаем, если канал done закрыт
//					}
//				}
//			}(i, stage, in)
//			in = channels[i] // Передаем выходной канал текущего стейджа как входной для следующего
//		}
//
//		// Передаем данные из последнего стейджа в выходной канал
//		go func() {
//			for v := range in {
//				select {
//				case <-done:
//					return
//				case out <- v:
//				}
//			}
//		}()
//	}()
//
//	return out
//}
