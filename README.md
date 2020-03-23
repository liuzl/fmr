# FMR: Functional Meaning Representation & Semantic Parsing Framework
[![GoDoc](https://godoc.org/github.com/liuzl/fmr?status.svg)](https://godoc.org/github.com/liuzl/fmr)[![Go Report Card](https://goreportcard.com/badge/github.com/liuzl/fmr)](https://goreportcard.com/report/github.com/liuzl/fmr)

## Projects that uses FMR

### mathsolver
* codes: https://github.com/liuzl/mathsolver
* demo: https://mathsolver.zliu.org/

## What is semantic parsing?
Semantic parsing is the process of mapping a natural language sentence into an intermediate logical form which is a formal representation of its meaning.

The formal representation should be a detailed representation of the complete meaning of the natural language sentence in a fully formal language that:

* Has a rich ontology of types, properties, and relations.
* Supports automated reasoning or execution.

## Representation languages
Early semantic parsers used highly domain-specific meaning representation languages, with later systems using more extensible languages like Prolog, lambda calculus, lambda dependancy-based compositional semantics (λ-DCS), SQL, Python, Java, and the Alexa Meaning Representation Language. Some work has used more exotic meaning representations, like query graphs or vector representations.

### FMR, a formal meaning representation language
* FMR stands for  functional meaning representation
* Context-Free Grammar for bridging NL and FMR
* *[VIM Syntax highlighting for FMR grammar file](https://github.com/liuzl/vim-fmr)*

## Tasks
* Grammar checkers
* Dialogue management
* Question answering
* Information extraction
* Machine translation

## What can FMR do, a glance overview
```js
// semantic parsing
"五与5.8的和的平方的1.5次方与two的和减去261.712" =>
nf.math.sub(
  nf.math.sum(
    nf.math.pow(
      nf.math.pow(
        nf.math.sum(
          5,
          nf.math.to_number("5.8")
        ),
        2
      ),
      nf.math.to_number("1.5")
    ),
    2
  ),
  nf.math.to_number("261.712")
); // denotation: 1000

// slot filling
"从上海到天津的机票" => nf.flight("上海", "天津");
"到重庆，明天，从北京" => nf.flight("北京", "重庆");
"到上海去" => nf.flight(null, "上海");
```

## References
* [Semantic Parsing: Past, Present, and Future](http://yoavartzi.com/sp14/slides/mooney.sp14.pdf), Raymond J. Mooney, 2014
* [Introduction to semantic parsing](https://github.com/liuzl/fmr-files/blob/master/cs224u-2019-intro-semparse.pdf), Bill MacCartney, 2019
* [Bringing machine learning and compositional semantics together](https://web.stanford.edu/~cgpotts/manuscripts/liang-potts-semantics.pdf), Percy Liang and Christopher Potts, 2014
* [SippyCup: A semantic parsing tutorial](https://github.com/wcmac/sippycup), Bill MacCartney, 2015
* [Semantic parsing in your browser](https://www.cs.toronto.edu/~muuo/writing/semantic-parsing-in-your-browser/), Muuo Wambua, 2018
