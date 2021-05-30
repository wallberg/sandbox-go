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
        nodes = list(g)
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

    # Second reduction: Backtrack to discover all simple paths, from
    # v, that contribute only distinct words. Remove v. Repeat until no nodes
    # remain. Keep track of max length cycles as we go.

    print()
    print("Second Reduction")

    words = [None] * (g.size() * 2)

    def unique_words(x, level, successor):
        """ Determine if all word pairs return unique words. """

        nonlocal g, words

        # Get words from previous level (already confirmed unique)
        if level > 2:
            keydict = g.get_edge_data(x[level-3], x[level-2])
            i = 2 * (level-3)
            words[i] = keydict['word1']
            words[i+1] = keydict['word2']

        # Get the new words
        keydict = g.get_edge_data(x[level-2], successor)
        word1, word2 = keydict['word1'], keydict['word2']

        # See if word1 and word2 are unique
        i = 2 * (level-2)
        for j in range(0, i):
            if words[j] == word1 or words[j] == word2:
                return False

        words[i], words[i+1] = word1, word2

        return True

    max_level = 0

    def S(n, level, x):

        """ Return values at level, which hold true for
        x_1 to x_level. """

        nonlocal g, max_level, candidate

        if level == 1:
            # print(candidate)
            return [candidate]

        values = []
        # print(x[0:level-1])

        for successor in g[x[level-2]]:

            # Does this transition contribute unique words?
            if unique_words(x, level, successor):

                if successor == candidate:
                    # Found path with unique words

                    # We are looking for cycles of maximum length
                    if level >= max_level:
                        max_level = level

                        print()
                        print("Words:", level-1)
                        i = 2 * (level - 1)
                        print("  ", words[0:i:2])
                        print("  ", words[1:i:2])

                else:
                    # Add it to the returned values
                    values.append(successor)

        return values

    # Find the longest cycle
    nodes = list(g)
    for candidate in nodes:
        # Search for a simple path, from candidate, that contributes
        # only distinct words
        for _ in walkers_backtrack(g.size(), S): pass

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
