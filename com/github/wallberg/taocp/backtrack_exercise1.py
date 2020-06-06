# -*- coding: utf-8 -*-
import unittest
from itertools import islice

from com.github.wallberg.taocp.Backtrack import basic_backtrack

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming
'''


def n_tuples(seq, n=None):
    ''' Generate all n-tuples using backtracking. '''

    def domain(n, k):
        return tuple(seq)

    def property(n, level, x):
        return True

    if n is None:
        n = len(seq)

    assert 0 < n and n <= len(seq)

    for x in basic_backtrack(n, domain, property):
        yield x


def permutations(seq, n=None):
    ''' Generate all permutations of distinct items using backtracking. '''

    def domain(n, k):
        return tuple(seq)

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''
        if level == 1:
            return True

        for k in range(1, level):
            if x[k-1] == x[level-1]:
                return False

        return True

    if n is None:
        n = len(seq)

    assert 0 < n and n <= len(seq)

    for x in basic_backtrack(n, domain, property):
        yield x


def combinations(seq, n=None):
    ''' Generate all combinations of size n using backtracking. '''

    def domain(n, k):
        return tuple(seq[:len(seq) + 1 - k])

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''
        if level == 1:
            return True

        if x[level-2] <= x[level-1]:
                return False

        return True

    if n is None:
        n = len(seq)

    assert 0 < n and n <= len(seq)

    for x in basic_backtrack(n, domain, property):
        yield x


def partitions(n):
    ''' Generate all integer partitions of n. '''

    def domain(n, k):
        return tuple(range(0, n // k + 1))

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''
        if level > 1 and x[level-2] < x[level-1]:
            return False

        total = sum(x[:level])
        if total > n:
            return False

        if n - (n - level) * x[level-1] > total:
            return False

        return True

    for x in basic_backtrack(n, domain, property):
        yield tuple(value for value in x if value > 0)


def set_partitions(seq):
    ''' Generate all set partitions of seq using restricted growth strings. '''

    def domain(n, k):
        return tuple(range(k))

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''
        if level > 1 and 1 + max(x[:level-1]) < x[level-1]:
            return False

        return True

    n = len(seq)

    for value in basic_backtrack(n, domain, property):
        result = []
        for x in range(n):
            t = tuple(seq[i] for i in range(n) if value[i] == x)
            if len(t) > 0:
                result.append(t)
        yield tuple(result)


def nested_parentheses(n):
    ''' Generate all nested parentheses with n pairs. '''

    def domain(n, k):
        return tuple(range(1, 2*k))

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''
        if level > 1 and x[level-1] <= x[level-2]:
            return False

        return True

    for value in basic_backtrack(n, domain, property):
        result = [')'] * (n * 2)
        for i in range(1, n+1):
            result[value[i-1]-1] = '('
        yield ''.join(result)


class Test(unittest.TestCase):

    def test_n_tuples(self):
        result = list(n_tuples((1, 2)))
        self.assertEqual(result, [(1, 1), (1, 2), (2, 1), (2, 2)])

        result = list(n_tuples((1, 2, 3)))
        self.assertEqual(len(result), 27)

        result = list(n_tuples((1, 2, 3), 2))
        self.assertEqual(result, [(1, 1), (1, 2), (1, 3), (2, 1), (2, 2),
                                  (2, 3), (3, 1), (3, 2), (3, 3)])

        result = list(n_tuples((1, 2, 3), 1))
        self.assertEqual(result, [(1, ), (2, ), (3, )])

    def test_permutations(self):
        result = list(permutations((1, 2, 3)))
        self.assertEqual(result, [(1, 2, 3), (1, 3, 2), (2, 1, 3), (2, 3, 1),
                                  (3, 1, 2), (3, 2, 1)])

        result = list(permutations((1, 2, 3, 4)))
        self.assertEqual(len(result), 24)

        result = list(permutations((1, 2, 3), 2))
        self.assertEqual(result, [(1, 2), (1, 3), (2, 1), (2, 3),
                                  (3, 1), (3, 2)])

    def test_combinations(self):
        result = list(combinations((1, 2, 3)))
        self.assertEqual(result, [(3, 2, 1)])

        result = list(combinations((1, 2, 3, 4), 2))
        self.assertEqual(result, [(2, 1), (3, 1), (3, 2), (4, 1), (4, 2),
                                  (4, 3)])

        result = list(combinations((1, 2, 3, 4, 5, 6), 3))
        self.assertEqual(len(result), 20)

    def test_partitions(self):
        result = list(partitions(4))
        self.assertEqual(result, [(1, 1, 1, 1), (2, 1, 1), (2, 2), (3, 1),
                                  (4, )])

        result = list(partitions(8))
        self.assertEqual(len(result), 22)

    def test_set_partitions(self):
        result = list(set_partitions(list(range(3))))
        self.assertEqual(result, [((0, 1, 2), ),
                                  ((0, 1), (2, )),
                                  ((0, 2), (1, )),
                                  ((0, ), (1, 2)),
                                  ((0, ), (1, ), (2, ))])

    def test_nested_parentheses(self):
        result = list(nested_parentheses(3))
        self.assertEqual(result, ['((()))', '(()())', '(())()',
                                  '()(())', '()()()'])

        result = list(nested_parentheses(4))
        self.assertEqual(len(result), 14)


if __name__ == '__main__':
    unittest.main(exit=False)
