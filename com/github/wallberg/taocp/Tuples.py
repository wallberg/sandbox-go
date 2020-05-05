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


if __name__ == '__main__':
    unittest.main(exit=False)
