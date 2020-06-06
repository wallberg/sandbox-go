# -*- coding: utf-8 -*-
import unittest
from itertools import islice

from com.github.wallberg.taocp.Trie import WordTrie, PrefixTrie

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming
'''


def basic_backtrack(n, domain, property):
    '''
    Algorithm B. Basic backtrack.

    n - sequence length we will be generating

    domain(n, k) - function which generates x_1 to x_k for D_k

    properties(n, level, k) - function to test if sequence x_1 to x_l
       holds as true
    '''

    # B1 [Initialize]
    level = 1
    x = [None] * n
    domains = [domain(n, k) for k in range(1, n+1)]
    domains_i = [0] * n

    goto = 'B2'
    while True:

        if goto == 'B2':
            # [Enter level l.]

            if level > n:
                yield tuple(x)
                goto = 'B5'
            else:
                x[level-1] = domains[level-1][domains_i[level-1]]
                goto = 'B3'

        elif goto == 'B3':
            # [Try x_l.]

            if property(n, level, x):
                # print('B3 property holds')
                level += 1
                goto = 'B2'
            else:
                goto = 'B4'

        elif goto == 'B4':
            # [Try again.]

            if domains_i[level-1] < len(domains[level-1]) - 1:
                domains_i[level-1] += 1
                x[level-1] = domains[level-1][domains_i[level-1]]
                goto = 'B3'
            else:
                goto = 'B5'

        elif goto == 'B5':
            # [Backtrack.]

            level -= 1

            if level > 0:
                if level < n:
                    x[level] = None
                    domains_i[level] = 0
                goto = 'B4'
            else:
                return


def n_queens(n):
    ''' Solve the n-queens problems using backtracking. '''

    def domain(n, k):
        ''' Domain is each value 1..n for every k '''
        return tuple(range(1, n+1))

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''
        if level == 1:
            return True

        for k in range(1, level):
            delta = abs(x[k-1] - x[level-1])
            if delta == 0:
                # Same column
                return False
            elif delta == level-k:
                # Same diagonal
                return False

        return True

    for x in basic_backtrack(n, domain, property):
        yield x


def walkers_backtrack(n, S, stats=None):
    '''
    Algorithm W. Walker's backtrack.

    n - sequence length we will be generating

    S(n, level, x) - function to return valid next values at level l, given
                     the existing values in x_1 .. x_(l-1)
    '''

    # W1. [Initialize.]
    level = 1
    x = [None] * n
    domains = [None] * n
    domains_i = [0] * n

    if stats is not None:
        stats['level_count'] = [0] * n

    goto = 'W2'
    while True:

        if goto == 'W2':
            # [Enter level l.]

            if level > n:
                yield tuple(x)
                goto = 'W4'
            else:
                domains[level-1] = S(n, level, x)
                domains_i[level-1] = 0
                goto = 'W3'

        elif goto == 'W3':
            # [Try to advance.]

            if domains_i[level-1] < len(domains[level-1]):
                x[level-1] = domains[level-1][domains_i[level-1]]

                if stats:
                    stats['level_count'][level-1] += 1

                level += 1
                goto = 'W2'
            else:
                goto = 'W4'

        elif goto == 'W4':
            # [Backtrack.]

            level -= 1

            if level > 0:
                domains_i[level-1] += 1
                goto = 'W3'
            else:
                return


def n_queens_2(n):
    ''' Solve the n-queens problems using Walker's Backtrack. '''

    def S(n, level, x):
        ''' Return values at level, which hold true for x_1 to x_(level-1) '''

        if level == 1:
            return domain

        result = []
        for d in domain:

            # Test if valid
            valid = True
            for k in range(0, level-1):
                if d == x[k]:
                    # Same column
                    valid = False
                    break

                if abs(x[k] - d) == level-1-k:
                    # Same diagonal
                    valid = False
                    break

            if valid:
                result.append(d)

        return result

    domain = tuple(range(1, n+1))

    for x in walkers_backtrack(n, S):
        yield x


def langford_pairs(n):
    '''
    Algorithm L. Langford pairs. '''

    # L1. [Initialize.]
    x = [0] * (2 * n)
    p = [k + 1 for k in range(0, n)] + [0]
    level = 1
    j = None
    y = [None] * (2 * n)
    # print(f'L1. {level=}, {x=}, {p=}, {j=}, {y=} ')
    goto = 'L2'

    while True:

        if goto == 'L2':
            # [Enter level l.]

            k = p[0]
            if k == 0:
                # print(f'L2. visit {x=}')
                yield tuple(x)
                goto = 'L5'
            else:
                j = 0
                while x[level-1] < 0:
                    level += 1
                # print(f'L2. {level=}, {x=}, {p=}, {j=}, {y=} ')
                goto = 'L3'

        elif goto == 'L3':
            # [Try x_l = k.]

            if level + k + 1 > 2 * n:
                goto = 'L5'
            elif x[level+k] == 0:
                x[level-1] = k
                x[level+k] = -1 * k
                y[level-1] = j
                p[j] = p[k]
                level += 1
                # print(f'L3. {level=}, {x=}, {p=}, {j=}, {y=} ')
                goto = 'L2'
            else:
                goto = 'L4'

        elif goto == 'L4':
            # [Try again.]

            j = k
            k = p[j]
            # print(f'L4. {level=}, {x=}, {p=}, {j=}, {y=} ')
            if k != 0:
                goto = 'L3'
            else:
                goto = 'L5'

        elif goto == 'L5':
            # [Backtrack.]

            level -= 1

            if level > 0:
                while x[level-1] < 0:
                    level -= 1
                k = x[level-1]
                x[level-1] = 0
                x[level+k] = 0
                j = y[level-1]
                p[j] = k
                # print(f'L5. {level=}, {x=}, {p=}, {j=}, {y=} ')
                goto = 'L4'
            else:
                return


class Test(unittest.TestCase):

    def test_n_queens(self):
        result = list(n_queens(4))
        self.assertEqual(result, [(2, 4, 1, 3), (3, 1, 4, 2)])

        result = list(n_queens(8))
        self.assertEqual(len(result), 92)

    def test_n_queens_2(self):
        result = list(n_queens_2(4))
        self.assertEqual(result, [(2, 4, 1, 3), (3, 1, 4, 2)])

        result = list(n_queens_2(8))
        self.assertEqual(len(result), 92)

    def test_long_n_queens(self):
        result = sum(1 for x in n_queens(12))
        self.assertEqual(result, 14200)

    def test_long_n_queens_2(self):
        result = sum(1 for x in n_queens_2(12))
        self.assertEqual(result, 14200)

    def test_langford_pairs(self):
        result = list(langford_pairs(3))
        self.assertEqual(result, [(2, 3, 1, -2, -1, -3),
                                  (3, 1, 2, -1, -3, -2)])

        result = sum(1 for x in langford_pairs(6))
        self.assertEqual(result, 0)

        result = sum(1 for x in langford_pairs(8))
        self.assertEqual(result, 150*2)

    def test_long_langford_pairs(self):
        result = sum(1 for x in langford_pairs(11))
        self.assertEqual(result, 17792*2)


if __name__ == '__main__':
    unittest.main(exit=False)
