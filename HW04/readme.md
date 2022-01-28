<p>Практические задания

<p>1. Реализовать генератор нагрузки, который нагружает базу типовыми запросами, подготовленными в рамках практического задания Урока 3.
<p>Ответ: Реализовал. Простенький запрос получился, но хоть как-то. 

<p>2. Используя генератор нагрузки, измерить производительность (QPS) с разными значениями для параметров MaxConns, MinConns pool'а соединений.
<p>Ответ: На 100 connect'ах база отвалилась. 

```
		cfg.MaxConns = 1
		cfg.MinConns = 1
			start attack
			duration: 10.288708042s
			threads: 1000
			queries: 30315
			QPS: 3031
		cfg.MaxConns = 10
		cfg.MinConns = 5
			start attack
			duration: 10.093315301s
			threads: 1000
			queries: 106583
			QPS: 10658
		cfg.MaxConns = 20
		cfg.MinConns = 10
			start attack
			duration: 10.084377254s
			threads: 1000
			queries: 117238
			QPS: 11723
		cfg.MaxConns = 40
		cfg.MinConns = 20
			start attack
			duration: 10.077240155s
			threads: 1000
			queries: 118083
			QPS: 11808
		cfg.MaxConns = 80
		cfg.MinConns = 40
			start attack
			duration: 10.104022697s
			threads: 1000
			queries: 106220
			QPS: 10622
		cfg.MaxConns = 100
		cfg.MinConns = 50
			start attack
			2022/01/28 18:14:16 failed to query data: failed to connect to `host=localhost user=techuser database=projects`: server error (FATAL: remaining connection slots are reserved for non-replication superuser connections (SQLSTATE 53300))
			exit status 1
```

<p>3. Подготовить отчет с информацией об используемом железе (используемый процессор, количество ядер, объём оперативной памяти, объём жесткого диска) и о пропускной способности сервера PostgreSQL для выбранного запроса. Представить отчет в файле throughput.txt.
<p>Ответ:

```
Не совсем понял третий пункт.
Параметры ПК i7 2600(4ядра 8 потоков 3400МГцб 16ГБ ОЗУ)
Контейнер postgres запускался на виртуальной машине которой выделено 4 ядра и 6ГБ ОЗУ. 
```