# -*- coding: utf-8 -*-
import unittest

'''
The Art of Computer Programming, Volume 4A, Combinatorial Algorithms,
Part 1, 2011

§7.2.1.1 Generating All n-Tuples
'''


def preprimes(m, n):
    '''
    Generate all preprime strings for an m-ary alphabet with n
    length tuples, along with the j index of the prime with n-extension.

    Algorithm F. Prime and preprime string generation.
    '''

    # F1. [Initialize.]
    a = [0] * (n + 1)
    a[0] = -1
    j = 1

    while True:
        # F2. [Visit.]
        yield(tuple(a[1:]), j)

        # F3. [Prepare to increase.]
        j = n
        while a[j] == m - 1:
            j -= 1

        # F4. [Add one.]
        if j == 0:
            return
        a[j] += 1

        # F5. [Make n-extension.]
        for k in range(j+1, n+1):
            a[k] = a[k-j]


def tuples_from_primes(m, n):
    '''
    Generate all m-ary, n-length tuples using Algorithm F. If we concatenate
    all prime strings whose j divides n, we get a de Bruijn cycle!
    '''

    cycle = []
    for prime, j in preprimes(m, n):
        if n % j == 0:
            cycle.extend(prime[0:j])

    cycle += cycle[0:n-1]

    for i in range(m**n):
        yield tuple(cycle[i:i+n])


class Test(unittest.TestCase):

    def test_preprimes(self):
        result = list(preprimes(2, 3))
        self.assertEqual(result, [((0, 0, 0), 1),
                                  ((0, 0, 1), 3),
                                  ((0, 1, 0), 2),
                                  ((0, 1, 1), 3),
                                  ((1, 1, 1), 1),
                                  ])

        result = list(preprimes(3, 4))
        self.assertEqual(len(result), 32)
        self.assertEqual(result[0], ((0, 0, 0, 0), 1))
        self.assertEqual(result[4], ((0, 0, 1, 1), 4))
        self.assertEqual(result[18], ((0, 2, 1, 0), 3))
        self.assertEqual(result[31], ((2, 2, 2, 2), 1))

    def test_tuples_from_primes(self):
        result = list(tuples_from_primes(2, 3))
        self.assertEqual(result, [(0, 0, 0),
                                  (0, 0, 1),
                                  (0, 1, 0),
                                  (1, 0, 1),
                                  (0, 1, 1),
                                  (1, 1, 1),
                                  (1, 1, 0),
                                  (1, 0, 0),
                                  ])

        result = list(tuples_from_primes(3, 4))
        self.assertEqual(len(result), 81)
        self.assertEqual(result[0], (0, 0, 0, 0))
        self.assertEqual(result[4], (1, 0, 0, 0))
        self.assertEqual(result[18], (0, 2, 1, 0))
        self.assertEqual(result[80], (2, 0, 0, 0))


if __name__ == '__main__':
    unittest.main(exit=False)
