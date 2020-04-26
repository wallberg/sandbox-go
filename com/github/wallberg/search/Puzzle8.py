# -*- coding: utf-8 -*-
import unittest

from com.github.wallberg.search import Node, Problem, Search

# 8-Puzzle problem


class State():

    def __init__(self, state=None):
        self.state = state

        if isinstance(state, list):
            # Flatten the 2-dimensional list and convert to tuple
            self.state = tuple([n for sl in state for n in sl])
        else:
            # Expect a tuple
            self.state = state

    def __lt__(self, other):
        return self.state < other.state

    def __eq__(self, other):
        return self.state == other.state

    def __hash__(self):
        return hash(self.state)

    def __str__(self):
        ls = [str(n) for n in self.state]
        ls.insert(6, ":")
        ls.insert(3, ":")
        return "".join(ls)


class StateTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.state1 = State()
        cls.state2 = State((1, 2, 3, 4, 5, 6, 7, 8, 0))
        cls.state3 = State((1, 2, 3, 4, 5, 6, 8, 7, 0))
        cls.state4 = State([[1, 2, 3], [4, 5, 6], [8, 7, 0]])

    def test_init(self):
        self.assertEqual(self.state1.state, None)
        self.assertEqual(self.state2.state, (1, 2, 3, 4, 5, 6, 7, 8, 0))
        self.assertEqual(self.state4.state, (1, 2, 3, 4, 5, 6, 8, 7, 0))

    def test_lt(self):
        self.assertTrue(self.state2 < self.state3)
        self.assertFalse(self.state3 < self.state4)

    def test_eq(self):
        self.assertTrue(self.state3 == self.state4)
        self.assertFalse(self.state2 == self.state3)

    def test_str(self):
        self.assertEqual(str(self.state2), "123:456:780")


class Puzzle8(Problem):

    def __init__(self, initial=None, goal=None):
        super().__init__(initial=initial)

        if goal is None:
            goal = State([[1, 2, 3], [4, 5, 6], [7, 8, 0]])

        self.goal = goal  # goal state for the search

    def is_goal(self, node):
        """ Determine if this is the goal node. """
        return node.state.state == self.goal.state

    def child_nodes(self, node):
        """ Generator for child nodes. """

        # Find the blank space
        for i in range(9):
            if node.state.state[i] == 0:
                row_blank = i // 3
                col_blank = i % 3
                break

        # Check for up to 4 ways to move the blank
        for operator in ["left", "right", "up", "down"]:

            # Determine the swap
            if operator == "left" and col_blank > 0:
                row_swap, col_swap = row_blank, col_blank - 1

            elif operator == "right" and col_blank < 2:
                row_swap, col_swap = row_blank, col_blank + 1

            elif operator == "up" and row_blank > 0:
                row_swap, col_swap = row_blank - 1, col_blank

            elif operator == "down" and row_blank < 2:
                row_swap, col_swap = row_blank + 1, col_blank

            else:
                # this move not possible
                continue

            value_swap = node.state.state[row_swap*3+col_swap]

            # Make the swap
            state_list = list(node.state.state)
            state_list[row_swap*3+col_swap] = 0
            state_list[row_blank*3+col_blank] = value_swap

            state_new = State(state=tuple(state_list))

            # Return the new child node
            yield Node(state=state_new,
                       parent=node,
                       operator=operator,
                       g=node.g + 1,
                       h=self.manhattan_distance(state_new))

    def manhattan_distance(self, state):
        h = 0

        # Loop through the state positions
        for row in range(3):
            for col in range(3):
                #  Skip the blank
                if state.state[row*3+col] != 0:
                    # Loop through the goal state positions
                    for rowg in range(3):
                        for colg in range(3):
                            # Check for a match
                            if state.state[row*3+col] \
                               == self.goal.state[rowg*3+colg]:
                                distance = abs(rowg - row) + abs(colg - col)
                                h += distance
                                # if (state.state[row][col] == 8
                                # and distance == 0):
                                #     h += 2
        return h


class Puzzle8Test(unittest.TestCase):

    def test_init(self):
        p = Puzzle8()
        self.assertEqual(p.initial, None)

        p = Puzzle8(initial=1)
        self.assertEqual(p.initial, 1)

    def test_is_goal(self):
        problem = Puzzle8()

        state = State((1, 2, 3, 4, 5, 6, 7, 8, 0))
        node = Node(state=state)
        self.assertTrue(problem.is_goal(node))

        state = State((1, 2, 3, 4, 5, 6, 7, 0, 8))
        node = Node(state=state)
        self.assertFalse(problem.is_goal(node))

    def test_child_nodes(self):
        state = State((1, 2, 3, 4, 5, 6, 7, 8, 0))
        node = Node(state=state)
        problem = Puzzle8()

        nodes = list(problem.child_nodes(node))
        self.assertEqual(len(nodes), 2)

        node = nodes[0]
        self.assertEqual(node.state.state, (1, 2, 3, 4, 5, 6, 7, 0, 8))
        self.assertEqual(node.g, 1)
        self.assertEqual(node.operator, "left")

        node = nodes[1]
        self.assertEqual(node.state.state, (1, 2, 3, 4, 5, 0, 7, 8, 6))
        self.assertEqual(node.g, 1)
        self.assertEqual(node.operator, "up")

        state = State((1, 2, 3, 4, 0, 5, 6, 7, 8))
        node = Node(state=state)

        nodes = list(problem.child_nodes(node))
        self.assertEqual(len(nodes), 4)

        self.assertEqual(nodes[0].state.state, (1, 2, 3, 0, 4, 5, 6, 7, 8))
        self.assertEqual(nodes[1].state.state, (1, 2, 3, 4, 5, 0, 6, 7, 8))
        self.assertEqual(nodes[2].state.state, (1, 0, 3, 4, 2, 5, 6, 7, 8))
        self.assertEqual(nodes[3].state.state, (1, 2, 3, 4, 7, 5, 6, 0, 8))

    def test_manhattan_distance(self):
        problem = Puzzle8()

        state = State([[1, 2, 3], [4, 5, 6], [7, 8, 0]])
        self.assertEqual(problem.manhattan_distance(state), 0)

        state = State([[1, 2, 3], [4, 5, 6], [8, 7, 0]])
        self.assertEqual(problem.manhattan_distance(state), 2)

        state = State([[7, 2, 8], [5, 4, 1], [6, 3, 0]])
        self.assertEqual(problem.manhattan_distance(state), 16)

    def test_search_a_star(self):
        # Setup the problem
        state = State([[4, 1, 3], [7, 2, 5], [0, 8, 6]])
        node = Node(state=state)
        problem = Puzzle8(initial=node)
        search = Search(problem=problem)

        final = search.a_star()
        self.assertIsInstance(final, Node)
        self.assertEqual(final.g, 6)
        self.assertEqual(final.f(), 6)

        final = search.a_star(max_depth=5)
        self.assertIsNone(final)

        final = search.a_star(max_cost=5)
        self.assertIsNone(final)

    def test_search_a_star_id(self):
        # Setup the problem
        state = State([[4, 1, 3], [7, 2, 5], [0, 8, 6]])
        node = Node(state=state)
        problem = Puzzle8(initial=node)
        search = Search(problem=problem)

        final = search.a_star_id()
        self.assertIsInstance(final, Node)
        self.assertEqual(final.g, 6)
        self.assertEqual(final.f(), 6)

        final = search.a_star_id(max_cost=5)
        self.assertIsNone(final)


if __name__ == '__main__':
    unittest.main(exit=False)

    # Setup the problem
    state = State([[7, 2, 8], [5, 4, 1], [6, 3, 0]])
    node = Node(state=state)
    problem = Puzzle8(initial=node)
    search = Search(problem=problem)

    # Execute the search
    final = search.a_star()

    # Present the results
    if final is None:
        print("Search failed")
    else:
        print("Success")
        print(final.path())
