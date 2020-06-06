# -*- coding: utf-8 -*-
import unittest

from com.github.wallberg.taocp.Backtrack import walkers_backtrack

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming
'''


def exercise_68(stats=None):
    ''' Solve Exercise 68 using Walker's Backtrack. '''

    R = 0
    L = 1
    D = 2
    U = 3

    def S(n, level, x):
        ''' Return values at level, which hold true for x_1 to x_(level-1) '''

        nonlocal pointsto, preset, dump

        i = level - 1

        # Build set of possible values
        known, unknown = count(level, x, i)
        known_count = len(known)

        domain = set(range(known_count, known_count+unknown+1))

        if i in preset:
            domain &= set([preset[i]])

        # Check if further constrained by past R/D values or future preset
        # values
        if domain:
            for j in pointsfrom[i]:
                if j < level:
                    value_j = x[j]
                elif j in preset:
                    value_j = preset[j]
                else:
                    # future unknown value, so nothing to check
                    continue

                # Here's what we know so far
                known_j, unknown_j = count(level, x, j)

                # Try our different domain values
                for value in list(domain):
                    known_test = known_j | set([value])
                    known_count_test = len(known_test)
                    unknown_test = unknown_j - 1

                    if value_j not in range(known_count_test,
                                            known_count_test+unknown_test+1):

                        # Violates constraint, remove from consideration
                        domain.remove(value)

        return tuple(domain)

    def count(level, x, i):
        ''' Return counts for a cell: (unique digits, empty cells). '''

        nonlocal pointsto, preset

        known = set()  # set of unique known values
        unknown = 0  # count of unknown values

        for j in pointsto[i]:
            if j < level - 1:
                known.add(x[j])
            elif j >= level and j in preset:
                known.add(preset[j])
            else:
                unknown += 1

        return((known, unknown))

    def dump(x):
        ''' Print a matrix representation of x. '''
        for row in range(0, 10):
            print(x[row*10:(row+1)*10])

    # preset values
    preset = {0: 3, 1: 1, 2: 4, 4: 1, 6: 5, 8: 9,
              12: 2, 13: 6, 19: 5,
              26: 3, 27: 5, 28: 8, 29: 9,
              38: 7,
              40: 9, 42: 3,
              57: 2, 59: 3,
              61: 8,
              70: 4, 71: 6, 72: 2, 73: 6,
              80: 4, 86: 3, 87: 3,
              91: 8, 93: 3, 95: 2, 97: 7, 98: 9, 99: 5}

    # direction of arrow (row, col)
    points = [D, L, R, L, D, L, D, D, D, D,
              R, D, L, D, D, R, D, L, R, D,
              R, D, R, R, D, D, R, D, L, L,
              R, R, D, U, D, R, L, D, L, L,
              R, R, D, D, D, D, L, L, D, L,
              R, R, R, L, D, R, D, R, D, D,
              R, R, D, U, D, D, U, L, U, L,
              U, R, D, U, D, R, D, L, D, U,
              U, R, R, R, D, L, L, L, L, U,
              R, R, U, L, U, U, L, L, U, U]

    # cells a given cell points to (index)
    pointsto = [None] * 100
    pointsfrom = [None] * 100
    for j in range(100):
        pointsfrom[j] = []

    for row in range(0, 10):
        for col in range(0, 10):
            i = row * 10 + col

            if points[i] == R:
                pointsto[i] = tuple(range(i+1, (row+1)*10, 1))
            elif points[i] == L:
                pointsto[i] = tuple(range(i-1, (row*10)-1, -1))
            elif points[i] == D:
                pointsto[i] = tuple(range(i+10, 100, 10))
            else:  # points[i] == U:
                pointsto[i] = tuple(range(i-10, -1, -10))

            for j in pointsto[i]:
                pointsfrom[j].append(i)

    for x in walkers_backtrack(100, S, stats=stats):
        yield x


class Test(unittest.TestCase):

    def test_exercise_68(self):
        result = list(exercise_68())
        self.assertEqual(len(result), 1)

        self.assertEqual(result[0], (3, 1, 4, 3, 1, 3, 5, 5, 9, 5,
                                     7, 4, 2, 6, 1, 3, 5, 7, 1, 5,
                                     9, 4, 7, 6, 1, 2, 3, 5, 8, 9,
                                     7, 6, 3, 2, 1, 3, 5, 4, 7, 7,
                                     9, 8, 3, 4, 1, 2, 6, 7, 5, 9,
                                     4, 4, 4, 1, 1, 3, 4, 2, 4, 3,
                                     9, 8, 3, 5, 1, 2, 4, 7, 6, 9,
                                     4, 6, 2, 6, 1, 3, 2, 5, 2, 4,
                                     4, 3, 3, 3, 1, 3, 3, 3, 3, 5,
                                     9, 8, 4, 3, 1, 2, 6, 7, 9, 5))


if __name__ == '__main__':
    unittest.main(exit=False)
