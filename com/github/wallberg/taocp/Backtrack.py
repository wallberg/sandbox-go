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


def walkers_backtrack(n, S):
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


def word_rectangles(m, n, m_trie=None, n_trie=None):
    '''
    Find m x n word rectangles using Basic Backtrack, using a prefix trie
    as in Exercise 24 and orthogonal lists as in Exercise 25.
    '''

    def domain(n, k):
        '''
        Domain for k=1 is every word in m_trie; for k > 1, D is dynamic based
        on a and h.
        '''
        return domains[k-1]

    def property(n, level, x):
        ''' Test if new sequence x[0] to x[level-1] holds true. '''

        # print(f'{level=}, {x=}')

        if level == n:
            return True

        # Build list of possible next words, D_(l+1)
        m_words = None

        # Iterate over positions in the n_word in x[level-1]
        for k in range(m):

            # Build list of possible "next letters" for at x_k

            # c is the index of the letter at x_l,k
            c = ord(x[level-1][k]) - 96

            # Set new value of a[level], next letter in the prefix string
            node = n_trie.trie[a[level-1][k]][c]
            a[level][k] = node

            if node == 0:
                # No next letters
                return False

            if w[k][node] == 0:

                # Find m_words which match at position k
                m_words_k = None

                # Iterate over next possible letters
                for c_h in range(1, 27):
                    if n_trie.trie[a[level][k]][c_h] != 0:

                        # Found possible next letter
                        # Check if we have m_words which match
                        if h[k][c_h-1]:

                            # Keep track of all words which match at each
                            # letter
                            if m_words_k is None:
                                m_words_k = set(h[k][c_h-1])
                            else:
                                m_words_k |= h[k][c_h-1]

                w[k][node] = set() if m_words_k is None else m_words_k

            if len(w[k][node]) == 0:
                # No matching m_words at k
                return False

            if m_words is None:
                m_words = set(w[k][node])
            else:
                m_words &= w[k][node]

            if len(m_words) == 0:
                return False

        # Set D_(l+1)
        domains[level].clear()
        domains[level].extend(m_words)
        domains[level].sort()

        return True

    # Load m length word trie
    if m_trie is None:
        m_trie = WordTrie()
        if m == 5:
            m_trie.load_sgb_words()
        else:
            m_trie.load_ospd4_words(m)
    else:
        assert isinstance(m_trie, WordTrie)

    # Load n length prefix trie
    if n_trie is None:
        n_trie = PrefixTrie()
        if n == 5:
            n_trie.load_sgb_words()
        else:
            n_trie.load_ospd4_words(n)
    else:
        assert isinstance(n_trie, PrefixTrie)

    # domain is all words in m_trie
    m_words = list(word for word in m_trie)

    domains = [0] * n
    domains[0] = m_words
    for level in range(1, n):
        domains[level] = []

    # a is an m x n lookup table for trie nodes corresponding to the prefixes
    # of the first l columns of partial solution x
    a = [[0] * m for j in range(n+1)]

    # h is an m x 26 lookup table for words which contain a specific letter
    # at position k
    h = [[0] * 26 for k in range(m)]
    for m_word in m_words:
        for k in range(m):
            c = ord(m_word[k]) - 96
            if not h[k][c-1]:
                h[k][c-1] = set()
            h[k][c-1].add(m_word)

    # w is a m x nodes(n_trie) lookup table for possible next m_words, given
    # the position k in an m_word and n_word prefix (= node in n_trie)
    w = [[0] * len(n_trie.trie) for k in range(m)]

    # for k in range(m):
    #     for c in range(1, 27):
    #         if h[k][c-1]:
    #             print(f'{k=},{c=},{h[k][c-1]=}')

    for x in basic_backtrack(n, domain, property):
        yield x


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

    def test_word_rectangles(self):
        m_trie = WordTrie()
        m_trie.add(['slums', 'total', 'awoke', 'tepid', 'using', 'faced'])

        n_trie = PrefixTrie()
        n_trie.add(['status', 'lowest', 'losers', 'lostly', 'utopia', 'making',
                    'sledge', 'facing'])

        result = list(word_rectangles(5, 6, m_trie, n_trie))
        self.assertEqual(len(result), 0)

        m_trie.add('stage')

        result = list(word_rectangles(5, 6, m_trie, n_trie))
        self.assertEqual(len(result), 1)
        self.assertEqual(result, [('slums', 'total', 'awoke', 'tepid', 'using',
                                  'stage')])

        result = list(word_rectangles(2, 2))
        self.assertEqual(len(result), 2177)

        # 5 x 6
        result = list(islice(word_rectangles(5, 6), 191))

        self.assertEqual(result[0],
                         ('aargh', 'blare', 'lapin', 'atilt', 'tense',
                          'edged'))

        self.assertEqual(result[190],
                         ('abaca', 'baths', 'bites', 'elude', 'sines',
                          'seers'))

        # result = list(word_rectangles(5, 6))
        # self.assertEqual(len(result), 625415)

        # for word in word_rectangles(5, 6):
        #     print(word)


if __name__ == '__main__':
    unittest.main(exit=False)
