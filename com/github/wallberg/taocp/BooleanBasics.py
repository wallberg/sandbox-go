# -*- coding: utf-8 -*-
import unittest

'''Explore Boolean Basics from The Art of Computer Programming, Volume 4a,
Combinatorial Algorithms, Part 1, 2011

§7.1.1 Boolean Basics
'''


def match_bit_pairs(v, j):
    '''Returns pairs of indexes into v for which their bits are all the same
    except at position j. Exercise 29.

    Arguments:
    * v - sequence of bitstrings (sort in ascending order)
    * j - bit which contains the single non-match
    '''

    # B1. [Initialize.]
    k, kp = 0, 0
    m = len(v)

    while True:

        # B2. [Find a zero.]
        while True:
            if k == m:
                return
            if v[k] & (1 << j) == 0:
                break
            k += 1

        # B3. [Make k-prime > k.]
        if kp <= k:
            kp = k + 1

        # B4. [Advance k-prime.]
        while True:
            if kp == m:
                return
            if v[kp] >= v[k] + (1 << j):
                break
            kp += 1

        # B5. [Skip past a big mismatch.]
        if v[k] ^ v[kp] >= 1 << (j+1):
            k = kp
            continue  # Goto B2

        # B6. [Record a match.]
        if v[kp] == v[k] + (1 << j):
            yield (k, kp)

        # B7. [Advance k.]
        k += 1


class Test(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.f22 = (0, 1, 4, 7, 12, 13, 14, 15)

    def test_match_bit_pairs(self):
        result = list(match_bit_pairs(self.f22, 0))
        self.assertEqual(result, [(0, 1), (4, 5), (6, 7)])

        result = list(match_bit_pairs(self.f22, 1))
        self.assertEqual(result, [(4, 6), (5, 7)])

        result = list(match_bit_pairs(self.f22, 2))
        self.assertEqual(result, [(0, 2)])

        result = list(match_bit_pairs(self.f22, 3))
        self.assertEqual(result, [(2, 4), (3, 7)])


if __name__ == '__main__':
    unittest.main(exit=False)
