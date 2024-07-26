// Package hypher enables creating AI agents as computational graphs.
//
// An agent is represented as a weighted Directed Acyclic Graph (DAG)
// which consists of nodes that perform a specific operation during
// agent execution.
//
// Given the agents are DAGs, they can form ensambles of agents
// through additional edges that link the agent DAGs as long
// as the resulting graph is also a DAG.
package hypher
