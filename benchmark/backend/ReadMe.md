<h1 align=center> Backend Benchmark, Specifications & POCs </h1>

## Introduction

The backend is one of the most important parts of a website. It handles all the internal and invisible management of a site. All the content, interaction, and main functionalities are managed by the backend.

As an important part, it is essential that it is advantageous to use, meaning it should bring positive aspects such as **performance**, **technical debt**, and **ease of integration**.

Here, we will test and compare these different technologies:
- Python (Flask)
- Golang (Mux)
- Javascript (NestJS)

## About

We will test a simple math project from Tek1, 101pong.

The goal is to have an API that does not just return a string, but is active and allows doing something useful.

It will repeat this operation an arbitrary number of times (here, 1,000,000) to make real performance comparisons.

It will then store the result in the user’s data based on the name given after the route `/hello/`.

Example: <br>
`/hello/paul?t0=[1,3,5]&t1=[7,9,-2]` <br>
Should return `Hello paul! your ball won't reach, computed in ...ms`

Another example: <br>
`/hello/paul?t0=[1.1,3,5]&t1=[-7,9,2]` <br>
Should return: `Hello paul! your incidence angle is 16.57, computed in ...ms`

It will first calculate the incidence angle, store the result in the database, and then display it.

The goal could be to compare:
- The performance of each framework's calculation
- The technical debt of each framework
- The speed of interaction with a database
- The simplicity of integration
and much more.

## Initialization

### Python

Installing the dependencies:
```sh
pip install -r requirements.txt
```

Running the application:
```sh
# in python/
python3 backend.py
```

### Golang

Running the application:
```sh
# in golang/
go run backend.go
```

### NestJS

Installing the dependencies:
```sh
npm i
```

Running the application:
```sh
# in nest/
npm run start
```

## Technical Debt

### Prelude

By "technical debt," I mean the ease of development and expansion of the
project in the future. I also refer to the flaws that could slow down the
development process due to the specifics of the language.

For example, verbosity represents **the amount of code needed to complete a
task**.

I also talk about maintainability. If adding a new feature is complicated,
the code is **less maintainable**, which causes delays in project development.
Code should be clear and not unnecessarily complicated.

This is, for instance, one of the reasons why assembly language is not a good
technology choice, as adding a new feature would take far too much time.

Verbosity is quite important, as it allows us, developers, to work quickly.

### Python

1 file, 53 lines of code

The Python code is the least verbose of all. Adding a feature with Flask is
quite simple, and the language is designed to be understandable by all
developers, even beginners.

However, the syntax makes some parts of the code less readable and confuses
some people. Moreover, since the language is not strictly typed, it can
sometimes be easy to make mistakes that won't appear until the API route is
tested.

### Golang

1 file, 73 lines of code

The Golang code is longer than Python's, but it's still quite easy to
understand. The fact that the language is very structured in terms of types
and syntax means the code is always clear, with fewer errors during program
execution.

Adding features seems simple. The code is clean and clear, and it’s not
difficult for a novice developer who has experience with other programming
languages to understand.

However, some functionalities are more verbose, like the error handling
conditions, compared to Python or JavaScript, where it's simpler to return the
string you want to display.

### NestJS

7 files, 126 lines of code

Creating a new route is very simple thanks to a terminal tool that allows you
to generate routes with a single command. Additionally, the code is easy to
understand, strictly typed, and TypeScript is simple for implementing new
functions.

However, the amount of code required to generate a simple project is completely
disastrous. A lot of code is necessary, which can make things more overwhelming
and confuse developers in a forest of code that only creates a simple API
route.

In the future, it could be conceivable that maintaining NestJS code with so
much boilerplate could become difficult.

## Performance

### Python

For the Python version using Flask, we tested the performance by making
1 000 000 GET requests that calculate the angle of incidence based on the
positions of a ball at two points (t0 and t1) in the XY plane.
The average response time was 1000ms.

Flask, while very simple to set up and use, shows its limitations in terms
of performance. As the number of requests increases, the time required to
process each one grows significantly. This is likely due to the dynamic
nature of Python and the overhead associated with interpreting the code at
runtime.

### Golang

The Golang implementation, on the other hand, performed much better. With an
average response time of just 23ms, it clearly outperforms both Python and
NestJS in terms of speed. Golang’s statically-typed nature and compiled runtime
allow it to handle large numbers of requests very efficiently, without the same
kind of performance penalties seen in interpreted languages like Python.

Golang also benefits from a simpler runtime and fewer abstractions,
which contributes to its overall faster response time. For our task of
calculating the angle of incidence 1 000 000 times, Golang stands out as a
highly efficient choice for backend development when speed is crucial.

### NestJS

NestJS, written in TypeScript and running on Node.js, had an average response
time of 42ms. While this is significantly faster than Flask, it still lags
behind Golang by a considerable margin. However, it does offer a good balance
between performance and developer experience, especially with TypeScript's
static typing and the extensive ecosystem of tools in the Node.js world.

NestJS's performance is adequate for many backend applications, though for
scenarios requiring rapid processing of a large number of requests, it may
not be the best choice compared to Golang.

## Conclusion

After analyzing the performance of each framework, it’s clear that **Golang**
is the superior choice in terms of raw speed, with an average response time of
only 23ms. **NestJS** is a good middle ground with reasonable performance
(42ms), and offers a great developer experience thanks to TypeScript and its
extensive tooling. However, it is still slower than Golang in high-performance
scenarios.

On the other hand, **Python** shows significant limitations in
performance, especially when handling a large number of requests (1000ms).
While Python is great for rapid development and prototyping, it doesn't fare
well when high-performance is a key requirement.

Thus, if **performance** is the primary concern, **Golang** should be the
go-to choice. For projects that value developer productivity, ease of use,
and maintainability over speed, **NestJS** and **Python** could
still be viable options depending on the specific needs of the project.
