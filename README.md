# go-hypher

[![Build Status](https://github.com/milosgajdos/go-hypher/workflows/CI/badge.svg)](https://github.com/milosgajdos/go-hypher/actions?query=workflow%3ACI)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/milosgajdos/go-hypher)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

> [!IMPORTANT]
> **THIS IS A WILD THEORETICAL EXPERIMENT**

`go-hypher` attempts to implement AI agents as computational graphs in Go.

> [!NOTE]
> The name of the project has been inspired by [Hypha](https://en.wikipedia.org/wiki/Hypha)
> which is understood to be a long, branching structure -- like a graph -- which enables
> communication in fungi by conducting electrical impulses through hyphae.

AI agent is represented as a weighted directed acyclic graph (DAG).

The graph nodes are the fundamental computational units: they execute an action anre return its result.
The action can be any computational task e.g. LLM prompt, function call, etc.

The actions are chained together through the graph edges. Each edge can have a weight assigned to it.

Nodes may have explicit input provided to them, though some nodes do not need any input to execute their action.

Some nodes are marked as input: they are the nodes where the agent execution starts.

There must be at least one output node: the node that generates the agent output i.e. the result
of its (node) actions. It's usually the leaf node in the graph.

When the agent runs, it first does a topological sort of its nodes and then executes each node action
in the order returned by the sort.

The topological sort is done in ascending order by the number of node incoming edges i.e. in-degrees of nodes.
Naturally, the nodes that have no incoming edges are executed first and usually have inputs provided, though,
the inputs may also come from the node Op(eration) they execute.

The output of each node is passed to all its successive nodes which in turn process it, along with their own
input, and then execute their own action and pass the output to the next node: the agent graphs works
in a similar way to a data stream -- the data flows from the inputs, are transformed by intermediate nodes
on their way to the output node. The output node provides the agent execution result.

Given each agent is a graph it should be possible to create ensemble of agents by composing agent graphs
to perform arbitrarily complex tasks.

Furthermore, in a similar way as the neural networks, we should be able to improve the agent operation
through output feedback using some sort of evaluation of output and some type of optimiser.

In the context of LLM agents, we could theoretically improve agent operation by tuning its node prompts
when running the agent in some simulation environments through several cycles: agents could maintain
their prompts history and eval result for each simulation run and the final tuned agent would then simply
pick the prompt that performs the best for the given task.
