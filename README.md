# Turbocycle
Это библиотека, предназначенная для быстрого расчета параметров цикла газотурбинных установок произвольного цикла.

## Основа архитектуры:
В основе архитектуры данной библиотеки лежит представление всех компонентов газотурбинной установки в виде узлов `Node`, связанных
между собой. Каждый узел имеет набор портов, часть которых он декларирует как входные, а часть - как выходные. 
Кроме этого каждый узел определяет внутренний алгоритм, согласно которому он обновляет состояние выходных портов в 
соответствии с состоянием входных портов.

Достоинством такой архитектуры является ее расширяемость: в силу высокой абстрактности интерфейса `Node`, библиотека 
доступных узлов может быть легко расширена за счет описания новых конструктивных (и не только) узлов ГТД. В принципе,
возможно также использование данной схемы в областях, не связанных с газотурбинной техникой.

## Примеры использования:
### Создание источника газа:
```go
var gasSource = source.NewComplexGasSourceNode(
  gases.GetAir(), // газ 
  288,            // температура газа 
  1e5             // давление газа
)
```

### Создание компрессорного узла
```go
var compressor = constructive.NewCompressorNode(
  0.86,     // КПД компрессора 
  6,        // степень повышения давления
  0.05      // точность расчета теплофизических параметров газа
)
```

### Соединение узлов
```go
var gasSource = source.NewComplexGasSourceNode(
  gases.GetAir(), // газ 
  288,            // температура газа 
  1e5,            // давление газа
)

var compressor = constructive.NewCompressorNode(
  0.86,     // КПД компрессора 
  6,        // степень повышения давления
  0.05,     // точность расчета теплофизических параметров газа
)

// Команда ниже соединяет комплексны газовый выход источника газа с комплексным газовым входом компрессора
// По комплексным газовым портам (в отличие от редко используемых обычных газовых портом) передается 4 параметра:
// сам газ, его температура, давление и относительный расход (по обычному газовому порту передается только газ)
core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())  
```

### Создание газотурбинной установки по двухвальной схеме со свободной турбиной.
```go
var nodes = make(map[string]core.Node)  // хранилище узлов (используется для соотнесения узлов с их названиями)

var gasSource1 = source.NewComplexGasSourceNode(gases.GetAir(), 300, 1e5) // источник газа на вход в компрессор
nodes["gasSource1"] = gasSource1

var gasSource2 = source.NewComplexGasSourceNode(gases.GetAir(), 300, 1e5) // источник газа на выход свободной турбине
nodes["gasSource2"] = gasSource2

var gasSink = sink.NewComplexGasSinkNode()  // сток газа на выход свободной турбины
nodes["gasSink"] = gasSink

var compressor = constructive.NewCompressorNode(0.86, 6, 0.05) // компрессорный узел 
nodes["compressor"] = compressor

var turbine = constructive.NewBlockedTurbineNode( // узел турбины компрессора
  0.92,       // КПД турбины 
  0.3,        // приведенная скорость на выходе (при расчете параметров цикла не используется, 
              // но нужна в случае определения статического давления на выходе)
  0.05,       // точность расчета теплофизических параметров газа
  func(constructive.TurbineNode) float64 {  // функция изменения расхода (возвращает разность между 
                                            // относительныv расходом на выходе и относительным расходом на входе)
	  return 0
  },
)
nodes["turbine"] = turbine

var burner = constructive.NewBurnerNode(  // узел камеры сгорания
  fuel.GetCH4(),  // топливо
  1400,           // температура торможения газа 
  300,            // температура топлива
  0.98,           // коэффициент сохранения полного давления 
  0.99,           // полнота сгорания 
  3,              // начальное значение коэффициента воздуха (используется как первое приближение в итеративных расчетах) 
  300,            // температура измерения теплоемкостей
  0.05,           // точность расчета теплофизических параметров газа
)
nodes["burner"] = burner

var freeTurbine = constructive.NewFreeTurbineNode(  // узел свободной турбины
  0.9,        // КПД турбины 
  0.3,        // приведенная скорость на выходе
  0.05,       // точность расчета теплофизических параметров
  func(constructive.TurbineNode) float64 {  // функция изменения расхода
	  return 0
  },
)
nodes["freeTurbine"] = freeTurbine

var powerSink1 = sink.NewPowerSinkNode()  // узел стока мощности
nodes["powerSink1"] = powerSink1

// особый узел, позволяющий собрать данные, идущие по раздельным каналам (канал давления, канал температуры, 
// канал газа, канал относительного расхода) в один комплексный газовы канал 
var assembler = helper.NewGasStateAssemblerNode()  
nodes["assembler"] =assembler

// особый узел, позволяющий разделить данные, идущие по комплексным газовым портам на 4 отдельных канала: 
// газ, давление, температура и относительный расход 
var disassembler = helper.NewGasStateDisassemblerNode()
nodes["disassembler"] = disassembler

var tSink = sink.NewTemperatureSinkNode() // узел стока температуры
nodes["tSink1"] = tSink
var mSink = sink.NewMassRateRelSinkNode() // узел стока относительного расхода
nodes["mSink"] = mSink
var gSink = sink.NewGasSinkNode()         // узел стока газа
nodes["gSink"] = gSink

var hub = helper.NewHubNode(        // разветвитель: подает данные на входе на оба своих выходе
  states.NewPressurePortState(1e5), // начальное состояние на выходах
)
nodes["hub"] = hub

// далее идет сборка установки
// соединить газовых вход компрессора с атмосферой
core.Link(compressor.ComplexGasInput(), gasSource1.ComplexGasOutput())
// соединить газовых выход компрессора с камерой сгорания
core.Link(compressor.ComplexGasOutput(), burner.ComplexGasInput())
// соединить газовых выход камеры сгорания с газовым входом турбины компрессора
core.Link(burner.ComplexGasOutput(), turbine.ComplexGasInput())
// соединить газовы выход турбины компрессора с газовым входом свободной турбины
core.Link(turbine.ComplexGasOutput(), freeTurbine.ComplexGasInput())
// соединить источник газа для свободной турбины со входом "разборщика"
core.Link(gasSource2.ComplexGasOutput(), disassembler.ComplexGasPort())
// замкнуть выходы разборщика по температуре, газу и относительному расходу на стоки, так как 
// свободной турбине не выходе требуется только давление
core.Link(disassembler.TemperaturePort(), tSink.TemperatureInput())
core.Link(disassembler.GasPort(), gSink.GasInput())
core.Link(disassembler.MassRateRelPort(), mSink.MassRateRelInput())
// поскольку давление на выходе из свободной турбины необходимо как в качестве граничного условия для расчета
// свободной турбины, так и в качестве результата расчета, подадим его на вход разветвителя
core.Link(disassembler.PressurePort(), hub.Inlet())
// подадим давление с первого выхода разветвителя на свободную турбину
core.Link(hub.Outlet1(), freeTurbine.PressureOutput())
// подадим давление со второго выхода разветвителя на вход "сборщика" по давлению
core.Link(hub.Outlet2(), assembler.PressurePort())
// соединим выходы свободной турбины по температуре, относительному расходу и газу с соответствующим входами сборщика
core.Link(freeTurbine.TemperatureOutput(), assembler.TemperaturePort())
core.Link(freeTurbine.MassRateRelOutput(), assembler.MassRateRelPort())
core.Link(freeTurbine.GasOutput(), assembler.GasPort())
// соединим выход сборщика с газовым стоком
core.Link(regenerator.ComplexGasPort(), gasSink.ComplexGasInput())

// соединим мощностные порты компрессора и турбины компрессора
core.Link(compressor.PowerOutput(), turbine.PowerInput())
// соединим мощностной выход свободной турбины с мощностным стоком
core.Link(freeTurbine.PowerOutput(), powerSink1.PowerInput())

// далее идет расчет цикла
// создать вычислительную сеть
var network = core.NewNetwork(nodes)

// решить 
// первый результат `converged` показывает, сошлось ли решение за максимальное количество итераций
// второй результат `err` - ошибка в случае, если исходная сеть была неправильно собрана:
// остались свободные порты или неразорванные петли (в данном примере петли отсутствуют)
var converged, err = network.Solve(
  0.1,      // коэффициент релаксации 
  100,      // максимальное число итераций 
  0.05      // точность расчета (насколько близки должны быть между собой последние итерации) 
)

// вывод решения в на экран (в формате json)
var b, _ = json.MarshalIndent(network, "", "    ")
	os.Stdout.Write(b)
```

### Пример результата решения (для другой схемы)
```json
{
    "assembler": {
        "ComplexGasPortState": {
            "p_stag": 102040.81632653062,
            "t_stag": 1157.7135639335263,
            "mass_rate_rel": 1.0188939094595912
        },
        "PressurePortState": {
            "p_stag": 102040.81632653062
        },
        "TemperaturePortState": {
            "t_stag": 1157.7135639335263
        },
        "GasPortState": {}
    },
    "breaker": {
        "State": {
            "p_stag": 102040.81632653062,
            "t_stag": 1157.7135639335263,
            "mass_rate_rel": 1.0188939094595912
        }
    },
    "burner": {
        "gas_input_state": {
            "p_stag": 800000,
            "t_stag": 1100.3121863924125,
            "mass_rate_rel": 1
        },
        "gas_output_state": {
            "p_stag": 792000,
            "t_stag": 1800,
            "mass_rate_rel": 1.0188939094595912
        },
        "alpha": 3.0620507120192766,
        "fuel_mass_rate_rel": 0.018893909459591178,
        "eta_burn": 0.99,
        "sigma": 0.99
    },
    "compressor": {
        "gas_input_state": {
            "p_stag": 100000,
            "t_stag": 300,
            "mass_rate_rel": 1
        },
        "gas_output_state": {
            "p_stag": 800000,
            "t_stag": 583.6997885223876,
            "mass_rate_rel": 1
        },
        "power_output_state": {
            "l_specific": -289677.6453939469
        },
        "eta_ad": 0.86,
        "pi_stag": 8,
        "mass_rate_rel": 1
    },
    "disassembler": {
        "ComplexGasPortState": {
            "p_stag": 102040.81632653062,
            "t_stag": 300,
            "mass_rate_rel": 1
        },
        "PressurePortState": {
            "p_stag": 102040.81632653062
        },
        "TemperaturePortState": {
            "t_stag": 300
        },
        "GasPortState": {}
    },
    "freeTurbine": {
        "gas_input_state": {
            "p_stag": 430395.38233249984,
            "t_stag": 1577.174969615828,
            "mass_rate_rel": 1.0188939094595912
        },
        "gas_output_state": {
            "p_stag": 102040.81632653062,
            "t_stag": 1157.7135639335263,
            "mass_rate_rel": 1.0188939094595912
        },
        "power_output_state": {
            "l_specific": 517981.4000426805
        },
        "pi_stag": 4.217874746858498,
        "l_specific": 517981.4000426805,
        "eta": 0.92
    },
    "gSink": {
        "gas_input_state": {}
    },
    "gasSink": {
        "gas_input_state": {
            "p_stag": 102040.81632653062,
            "t_stag": 666.2670958993333,
            "mass_rate_rel": 1.0188939094595912
        }
    },
    "gasSource1": {
        "gas_output_state": {
            "p_stag": 100000,
            "t_stag": 300,
            "mass_rate_rel": 1
        }
    },
    "gasSource2": {
        "gas_output_state": {
            "p_stag": 100000,
            "t_stag": 300,
            "mass_rate_rel": 1
        }
    },
    "hub": {
        "State": {
            "p_stag": 102040.81632653062
        }
    },
    "loss": {
        "gas_input_state": {
            "p_stag": 439178.9615637754,
            "t_stag": 1577.174969615828,
            "mass_rate_rel": 1.0188939094595912
        },
        "gas_output_state": {
            "p_stag": 430395.38233249984,
            "t_stag": 1577.174969615828,
            "mass_rate_rel": 1.0188939094595912
        },
        "sigma": 0.98
    },
    "mSink": {
        "mass_rate_input_state": {
            "mass_rate_rel": 1
        }
    },
    "outflow": {
        "gas_input_state": {
            "p_stag": 102040.81632653062,
            "t_stag": 300,
            "mass_rate_rel": 1
        },
        "gas_output_state": {
            "p_stag": 100000,
            "t_stag": 300,
            "mass_rate_rel": 1
        },
        "sigma": 0.98
    },
    "powerSink1": {
        "power_input_state": {
            "l_specific": 289677.6453939469
        }
    },
    "powerSink2": {
        "power_input_state": {
            "l_specific": 517981.4000426805
        }
    },
    "regenerator": {
        "hot_input_state": {
            "p_stag": 102040.81632653062,
            "t_stag": 1157.7135639335263,
            "mass_rate_rel": 1.0188939094595912
        },
        "cold_input_state": {
            "p_stag": 800000,
            "t_stag": 583.6997885223876,
            "mass_rate_rel": 1
        },
        "hot_output_state": {
            "p_stag": 102040.81632653062,
            "t_stag": 666.2670958993333,
            "mass_rate_rel": 1.0188939094595912
        },
        "cold_output_state": {
            "p_stag": 800000,
            "t_stag": 1100.3121863924125,
            "mass_rate_rel": 1
        },
        "sigma": 0.9
    },
    "tSink1": {
        "temperature_input_state": {
            "t_stag": 300
        }
    },
    "turbine": {
        "gas_input_state": {
            "p_stag": 792000,
            "t_stag": 1800,
            "mass_rate_rel": 1.0188939094595912
        },
        "gas_output_state": {
            "p_stag": 439178.9615637754,
            "t_stag": 1577.174969615828,
            "mass_rate_rel": 1.0188939094595912
        },
        "power_input_state": {
            "l_specific": -289677.6453939469
        },
        "power_output_state": {
            "l_specific": 289677.6453939469
        },
        "l_specific": 289677.6453939469,
        "eta": 0.92
    }
}
```
