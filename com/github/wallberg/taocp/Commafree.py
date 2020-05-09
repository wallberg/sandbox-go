# -*- coding: utf-8 -*-
import unittest

from itertools import product

from com.github.wallberg.taocp.Backtrack import walkers_backtrack

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming - Commafree codes
'''


def exercise34(words):
    ''' Find the largest commafree subset. '''

    def S(n, level, x):
        ''' Return values at level, which hold true for x_1 to x_(level-1) '''

        if level == 1:
            return domain

        elif x[level-2] is None:
            return [None]

        # Check all words after x_(l-1) for commafreeness
        result = []
        i = domain.index(x[level-2]) + 1
        while i < n:
            # Check if domain[i] is commafree with respect to x_1 .. x_(l-1)
            if exercise35(domain[i], x[:level-1]):
                result.append(domain[i])
            i += 1

        if len(result) > 0:
            return result
        else:
            return [None]

    domain = tuple(sorted(words))
    n = len(domain)

    max_subset = []
    for cf in walkers_backtrack(n, S):
        size = cf.index(None)
        if size > len(max_subset):
            max_subset = cf[:size]

    return max_subset


def exercise35(word, words):
    ''' Test if word is commafree with respect to the accepted words '''

    test_words = set(words + [word])
    for test1, test2 in product(test_words, test_words):
        testseq = test1 + test2

        for start in [0, 1]:
            count = 0
            for i in range(0, 4):
                if testseq[start+i:start+i+4] in test_words:
                    count += 1
            if count > 1:
                return False

    return True


class Test(unittest.TestCase):

    def test_exercise34(self):

        words = ['aced', 'babe', 'bade', 'bead', 'beef', 'cafe', 'cede',
                 'dada', 'dead', 'deaf', 'face', 'fade', 'feed']

        self.assertEqual(exercise34(words),
                         ('aced', 'babe', 'bade', 'bead', 'beef', 'cafe',
                          'cede', 'dead', 'deaf', 'fade', 'feed'))

    def test_exercise35(self):

        words = ['aced', 'babe', 'bade', 'bead', 'beef', 'cafe', 'cede',
                 'dada', 'dead', 'deaf', 'face', 'feed']

        self.assertTrue(exercise35('abcd', []))
        self.assertTrue(exercise35('abcd', ['efgh']))

        self.assertFalse(exercise35('cefa', words))
        self.assertFalse(exercise35('cece', ['cece']))
        self.assertFalse(exercise35('zzzz', []))


if __name__ == '__main__':
    unittest.main(exit=False)
