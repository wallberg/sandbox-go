# -*- coding: utf-8 -*-
import unittest

'''Explore Boolean Basics from The Art of Computer Programming, Volume 4a,
Combinatorial Algorithms, Part 1, 2011

§7.1.1 Boolean Basics
'''


def match_bit_pairs(v, j, start=0, stop=None):
    '''Returns pairs of indexes into v for which their bits are all the same
    except at position j. Exercise 29.

    Arguments:
    * v - integer sequence representing bitstrings (sorted in ascending order)
    * j - bit which contains the single non-match
    '''

    # B1. [Initialize.]
    m = len(v)

    if stop is None:
        stop = m

    k, kp = start, start

    while True:

        # B2. [Find a zero.]
        while True:
            if k == stop:
                return
            if v[k] & (1 << j) == 0:
                break
            k += 1

        # B3. [Make k-prime > k.]
        if kp <= k:
            kp = k + 1

        # B4. [Advance k-prime.]
        while True:
            if kp == stop:
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


def maximal_subcubes(n, v):
    '''Return all maximal subcubes (aka prime implicant) of v.
    Subcubes are represented as tuples (a, b) where a records the
    position of the asterisks and b records the bits in non-* positions.
    Exercise 30.

    Arguments:
    * v - integer sequence representing bitstrings (sorted in ascending order)
    * n - size in bits of bitstrings in v
    '''

    # P1. [Initialize.]
    m = len(v)

    # The current value of a being processed
    A = 0

    # Stack S contains |A| + 1 lists of subcubes; each list contains subcubes
    # with the same a value, in increasing order of a. This includes all lists
    # with wildcards plus the first list with a=0, equivalent to the input
    # v list of bitstrings
    S = [0] * (2*m + n)

    # Tag bits indicating a matching j-buddy, for each corresponding subcube
    # in S
    T = [0] * (2*m + n)

    # Determine the j-buddy pairs for the initial subcube list
    for j in range(n):
        for k, kp in match_bit_pairs(v, j):
            T[k] |= (1 << j)
            T[kp] |= (1 << j)

    # For each subcube, either output it as maximal or advance it
    # (with j-buddy) to the next subcube list, with additional wildcard
    s, t = 0, 0
    while s < m:
        if T[s] == 0:
            yield (0, v[s])
        else:
            S[t] = v[s]
            T[t] = T[s]
            t += 1
        s += 1
    S[t] = 0

    while True:
        # P2. [Advance A.]
        j = 0
        if S[t] == t:  # the topmost list is empty
            while j < n and A & (1 << j) == 0:
                j += 1
        while j < n and A & (1 << j) != 0:
            t = S[t] - 1
            A -= (1 << j)
            j += 1
        if j >= n:
            return
        A += (1 << j)

        # P3. [Generate list A.]
        r, s = t, S[t]
        for k, kp in match_bit_pairs(S, j, start=s, stop=r):
            x = (T[k] & T[kp]) - (1 << j)
            if x == 0:
                yield(A, S[k])
            else:
                t += 1
                S[t] = S[k]
                T[t] = x

        t += 1
        S[t] = r + 1


def str_subcube(n, a, b):
    '''Return string representation of an n-bit subcube in (a, b) form. '''

    s = ''
    for j in range(n-1, -1, -1):
        if a & (1 << j) == (1 << j):
            s += '*'
        else:
            s += '0' if b & (1 << j) == 0 else '1'
    return s


class Test(unittest.TestCase):

    F22 = (0, 1, 4, 7, 12, 13, 14, 15)  # "random" function 7.1.1-(22)
    F22_N = 4
    F22_SUBCUBES = [
        (1, 0),   # 000*
        (3, 12),  # 11**
        (4, 0),   # 0*00
        (8, 4),   # *100
        (8, 7)    # *111
    ]

    def test_match_bit_pairs(self):
        result = list(match_bit_pairs(self.F22, 0))
        self.assertEqual(result, [(0, 1), (4, 5), (6, 7)])

        result = list(match_bit_pairs(self.F22, 1))
        self.assertEqual(result, [(4, 6), (5, 7)])

        result = list(match_bit_pairs(self.F22, 2))
        self.assertEqual(result, [(0, 2)])

        result = list(match_bit_pairs(self.F22, 3))
        self.assertEqual(result, [(2, 4), (3, 7)])

    def test_maximal_subcubes(self):

        result = list(maximal_subcubes(4, (1, 2, 4, 8)))
        self.assertEqual(result, [(0, 1), (0, 2), (0, 4), (0, 8)])

        result = list(maximal_subcubes(self.F22_N, self.F22))
        self.assertEqual(result, self.F22_SUBCUBES)

        result = list(maximal_subcubes(5, list(range(32))))
        self.assertEqual(result, [(31, 0)])


if __name__ == '__main__':
    unittest.main(exit=False)
