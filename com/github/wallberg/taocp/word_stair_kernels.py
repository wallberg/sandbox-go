# -*- coding: utf-8 -*-
import logging
import sys
import re
import argparse

import networkx as nx

from com.github.wallberg.taocp.Backtrack import walkers_backtrack

"""
Explore Word Stair Kernels from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2.1 Dancing Links, Exercise 91
"""

logger = logging.getLogger('com.github.wallberg.taocp.word_stair_kernels')


def exercise_91(args):
    """ Find the longest cycle in a multi-set of right word stair kernels. """

    g = nx.DiGraph()

    kernels = set()

    # Read in each XCC provided kernel
    for line in sys.stdin:

        # Extract the kernel
        kernel = [0] * 14
        for option in line.strip().split(", "):
            option = option[1:-1].split(' ')

            if option[0] == "c1c2c6c10c13":
                kernel[0] = option[1][3]

            elif option[0] == "x3x4x5c2c3":
                kernel[1] = option[2][3]
                kernel[2] = option[3][3]

            elif option[0] == "c4c5c6c7c8":
                kernel[3] = option[1][3]
                kernel[4] = option[2][3]
                kernel[5] = option[3][3]
                kernel[6] = option[4][3]
                kernel[7] = option[5][3]

            elif option[0] == "c9c10c11c12x6":
                kernel[8] = option[1][3]
                kernel[9] = option[2][4]
                kernel[10] = option[3][4]
                kernel[11] = option[4][4]

            elif option[0] == "c13c14x7x8x9":
                kernel[12] = option[1][4]
                kernel[13] = option[2][4]

        kernel = "".join(kernel)

        # Add to set of distinct kernels
        if kernel not in kernels:

            # Add to the digraph of 9-tuple transitions, with the kernel
            # being the arc
            n1 = ''.join(kernel[0:7] + kernel[8:10])
            n2 = ''.join(kernel[2:3] + kernel[6:14])

            if args.right:
                # Right word stair
                word1 = kernel[3:8]
            else:
                # Left word stair
                word1 = kernel[7:2:-1]

            word2 = kernel[0:2] + kernel[5] + kernel[9] + kernel[12]

            g.add_edge(n1, n2, kernel=kernel, word1=word1, word2=word2)

            kernels.add(kernel)

    print("Kernels: ", len(kernels))

    def print_graph(g, verbose=False):
        print('nodes: ', len(g))
        print('edges: ', sum([1 for e in g.edges()]))
        if verbose:
            for n1, n2, kernel in g.edges(data='kernel'):
                print(f"{n1} -> {n2} {kernel}")

    print_graph(g, verbose=False)

    # First reduction: Get the largest induced subgraph for which every v has
    # positive in-degree and out-degree
    print()
    print("First Reduction")

    while True:
        changed = False
        nodes = list(n for n in g)
        for n in nodes:
            in_degree = sum([1 for s in g.predecessors(n)])
            out_degree = sum([1 for p in g.successors(n)])
            if in_degree == 0 or out_degree == 0:
                g.remove_node(n)
                changed = True

        if not changed:
            break

    if len(g) > 0:
        print_graph(g, verbose=False)

    # Second reduction: Find a node v of out-degree 1. Backtrack to discover a
    # simple path, from v, that contributes only distinct words. If there is no
    # such path, remove v from the graph and reduce it again.

    print()
    print("Second Reduction")

    def unique_words(pairlist):
        """ Determine if all word pairs return unique words. """

        nonlocal g

        words = set()

        for i in range(0, len(pairlist)-1):
            node1, node2 = pairlist[i], pairlist[i+1]

            keydict = g.get_edge_data(node1, node2)
            word1, word2 = keydict['word1'], keydict['word2']

            if word1 in words:
               return False
            words.add(word1)

            if word2 in words:
               return False
            words.add(word2)

        return True

    # Setup backtrack S() function
    class PathFound(Exception):
        pass

    max_cycle = []

    def get_s(reduction=True):
        """ Generate the S function for walkers_backtrack. """

        def S(n, level, x):

            """ Return values at level, which hold true for
            x_1 to x_(level-1). """

            nonlocal g, max_cycle

            if level == 1:
                return [candidate]

            values = []

            for successor in g[x[level-2]]:

                # Does this transition contribute unique words?
                path = x[0:level-1] + [successor]
                if unique_words(path):

                    if successor == candidate:
                        # Found path with unique words

                        if reduction:
                            # If reducing we want to stop here and indicate
                            # we found a path
                            raise PathFound()

                        else:
                            # We are looking for cycles of maximum length
                            cycle = path
                            if len(cycle) >= len(max_cycle):
                                max_cycle = cycle

                                # Gather the words
                                words = []
                                for i in range(0, len(cycle)-1):
                                    keydict = g.get_edge_data(cycle[i], cycle[i+1])
                                    words.append(keydict['word1'] + ":" + keydict['word2'])

                                print()
                                print("Cycle:", cycle)
                                print("Words:", words)
                    else:
                        # Add it to the returned values
                        values.append(successor)

            return values

        return S

    # # Loop until no node is removed
    # while True:

    #     node_removed = False

    #     # Iterate over all nodes
    #     for candidate in g:
    #         out_degree = sum([1 for p in g.successors(candidate)])
    #         if out_degree == 1:
    #             try:
    #                 # Search for a simple path, from candidate, that contributes
    #                 # only distinct words
    #                 for _ in walkers_backtrack(g.size(), get_s()):
    #                     pass

    #                 # No path found
    #                 g.remove_node(candidate)
    #                 node_removed = True
    #                 break

    #             except PathFound:
    #                 # Path found, so we can continue to the next candidate
    #                 pass

    #     if not node_removed:
    #         break

    # if len(g) > 0:
    #     print_graph(g, verbose=False)

    # Find the longest cycle
    nodes = list(g)
    for candidate in nodes:
        # Search for a simple path, from candidate, that contributes
        # only distinct words
        for _ in walkers_backtrack(g.size(), get_s(reduction=False)):
            pass

        # Now look for cycles which don't include this node
        g.remove_node(candidate)


if __name__ == '__main__':

    logger.addHandler(logging.StreamHandler())
    logger.setLevel(logging.INFO)

    # Setup command line arguments
    parser = argparse.ArgumentParser()

    parser.add_argument("-r", "--right",
                        action='store_true',
                        help="Process a right word stair; (default: a left word stair)")

    # Process command line arguments
    args = parser.parse_args()

    exercise_91(args)
