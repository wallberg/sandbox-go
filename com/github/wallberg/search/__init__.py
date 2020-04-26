# -*- coding: utf-8 -*-
import unittest
import string
import logging

from sortedcontainers import SortedList

# Search

logger = logging.getLogger('com.github.wallberg.search')


class Node():
    """ Node in the search tree """

    def __init__(self, state=None, parent=None, operator="noop", g=0, h=0):
        self.state = state  # state of the node, provide by Problem
        self.parent = parent  # parent of this node in the search tree
        self.operator = operator  # operator used to arrive at this state
        self.g = g  # g(n), path cost
        self.h = h  # h(n), heuristic (estimated future cost)

    def f(self):
        """ Total cost """
        return self.g + self.h

    def __lt__(self, other):
        return self.state < other.state

    def __eq__(self, other):
        return self.state == other.state

    def __hash__(self):
        return hash(self.state)

    def __repr__(self):
        return str(self)

    def __str__(self):
        return f"node: g={self.g} h={self.h}, f={self.f()}," + \
            f" op={self.operator}, state={self.state}"

    def path(self):
        """ Generate a string representation of the path from the parent
        node. """

        node = self
        nodes = [node]
        while node.parent is not None:
            node = node.parent
            nodes.insert(0, node)

        return "\n".join([str(node) for node in nodes])


class NodeTest(unittest.TestCase):

    def test_init(self):
        n = Node()
        self.assertEqual(n.state, None)
        self.assertEqual(n.parent, None)
        self.assertEqual(n.operator, "noop")
        self.assertEqual(n.g, 0)
        self.assertEqual(n.h, 0)

        n = Node(state=1, parent=2, operator=3, g=4, h=5)
        self.assertEqual(n.state, 1)
        self.assertEqual(n.parent, 2)
        self.assertEqual(n.operator, 3)
        self.assertEqual(n.g, 4)
        self.assertEqual(n.h, 5)

    def test_str(self):
        n = Node(state=1, parent=2, operator=3, g=4, h=5)
        self.assertEqual(str(n), "node: g=4 h=5, f=9, op=3, state=1")

    def test_path(self):
        n = Node()
        self.assertIsInstance(n.path(), str)


class Problem():
    """ Problem definition, to be overridden by the implementation. """

    def __init__(self, initial=None):
        self.initial = initial  # initial node in the search tree


class ProblemTest(unittest.TestCase):

    def test_initial(self):
        p = Problem()
        self.assertEqual(p.initial, None)

        p = Problem(initial=1)
        self.assertEqual(p.initial, 1)


class Search():
    """ Search routines to search based on the given Problem. """

    def __init__(self, problem=None):
        self.problem = problem  # definition of the problem to solve

    def general(self, max_depth=0, max_cost=0, key=None):
        """ General search. """

        # Setup the search
        self.queue = SortedList([self.problem.initial], key=key)
        self.allnodes = set([self.problem.initial])

        # Search
        self.count = 0
        while len(self.queue) > 0:

            # Get the next node to check
            node = self.queue.pop(0)
            self.count += 1

            if logger.isEnabledFor(logging.DEBUG):
                logger.debug("Next node is %s", node)

            # Check if this is the goal node
            if self.problem.is_goal(node):
                return node

            # Expand nodes
            for child in self.problem.child_nodes(node):

                # Don't add to the queue if we have seen this node already
                if child in self.allnodes:
                    continue

                # Don't add to the queue if max depth is exceeded
                if max_depth > 0 and child.g > max_depth:
                    continue

                # Don't add to the queue if max cost is exceeded
                if max_cost > 0 and child.f() > max_cost:
                    continue

                # Add to the queue
                self.queue.add(child)
                self.allnodes.add(child)

        # Goal not found
        return None

    # TODO: add searches for
    #   depth_first() - depth first
    #   depth_first_id() - depth first, with iterative deepening
    #   breadth_first() - breadth first

    def a_star(self, *args, **kwargs):
        """ A* search """

        kwargs['key'] = lambda node: node.f()

        return self.general(*args, **kwargs)

    def a_star_id(self, *args, **kwargs):
        """ A* search, with iterative deepening """

        max_cost_id = 0 if 'max_cost' not in kwargs else kwargs['max_cost']

        # Iteratively increase the maximum cost for an A* search
        cost = 0
        while True:
            cost += 1
            if max_cost_id > 0 and cost > max_cost_id:
                # Search failed
                return None

            # Execute the search
            kwargs['max_cost'] = cost
            final = self.a_star(*args, **kwargs)

            if final is not None:
                return final


class SearchTest(unittest.TestCase):

    def test_initial(self):
        s = Search()
        self.assertEqual(s.problem, None)

        s = Search(problem=1)
        self.assertEqual(s.problem, 1)


if __name__ == '__main__':
    unittest.main(exit=False)
