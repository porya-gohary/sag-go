<h1 align="center">
  <br>
  <a href="https://postimg.cc/0z1zFYn2"><img src="https://i.postimg.cc/kgBWVwt8/logo.png" alt="go-sag" width="300"></a>
  <br>
  Schedule-Abstraction Graph in GO
  <br>
</h1>

<h4 align="center">Unofficial implementation of schedule-abstraction graph using GO lang</h4>

<p align="center">
  <a href="https://github.com/porya-gohary/Multi-rate-DAG-Framework/blob/master/LICENSE.md">
    <img src="https://img.shields.io/hexpm/l/apa"
         alt="Gitter">
  </a>
    <img src="https://img.shields.io/badge/Made%20with-GO-orange">

</p>
<p align="center">
  <a href="#-required-packages">Dependencies</a> •
  <a href="#-build-instructions">Build</a> •
  <a href="#-input-format">Input Format</a> •
  <a href="#%EF%B8%8F-usage">Usage</a> •
  <a href="#-features">Features</a> •
  <a href="#-limitations">Limitations</a> •
  <a href="#-license">License</a>
</p>
<h4 align="center">NOTICE: THIS PROGRAM IS UNDER DEVELOPMENT...</h4>


## 📦 Required Packages
This program uses the following packages:

```
github.com/docopt/docopt-go
github.com/lfkeitel/verbose
gopkg.in/yaml.v3
```


## 📋 Build Instructions
After installing the required packages, simply build and run the code with:
```
go run ./main.go -j <input-file> 
```

## 📄 Input Format
This tool works with old SAG input format with csv format ([Example](./example/example3.csv)) and also new SAG input format with yaml format ([Example](./example/example3.yaml)).
Each input file describes a set of jobs. Each job is described by the following fields:
1.   **Task ID** — an arbitrary numeric ID to identify the task to which a job belongs
2.   **Job ID** — a unique numeric ID that identifies the job
3.   **Release min** — the earliest-possible release time of the job (equivalently, this is the arrival time of the job)
4.   **Release max** — the latest-possible release time of the job (equivalently, this is the arrival time plus maximum jitter of the job)
5.   **Cost min** — the best-case execution time of the job (can be zero)
6.   **Cost max** — the worst-case execution time of the job
7.   **Deadline** — the absolute deadline of the job
8.   **Priority** — the priority of the job (EDF: set it equal to the deadline)

## ⚙️ Usage


## 🔧 Features
- Classic single processor SAG.
- Single processor SAG with partial-order reduction.

## 🚧 Limitations
- For now, the framework just supports single processor.

## 📝 TODO
- [x] Implementation of uni-processor
- [x] Implementation of uni-processor with partial-order reduction
- [ ] Implement IIP
- [ ] Implement dependency
- [ ] Implement global multi-processor

## 📜 License
Copyright © 2022 [Pourya Gohari](https://pourya-gohari.ir)

This project is licensed under the Apache License 2.0 - see the [LICENSE.md](LICENSE.md) file for details