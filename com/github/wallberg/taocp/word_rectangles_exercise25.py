# -*- coding: utf-8 -*-
import unittest
from itertools import islice

from com.github.wallberg.taocp.Trie import WordTrie, PrefixTrie
from com.github.wallberg.taocp.Backtrack import basic_backtrack

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming
'''


def exercise25(m, n, m_trie=None, n_trie=None):
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

    def test_exercise25(self):
        m_trie = WordTrie()
        m_trie.add(['slums', 'total', 'awoke', 'tepid', 'using', 'faced'])

        n_trie = PrefixTrie()
        n_trie.add(['status', 'lowest', 'losers', 'lostly', 'utopia', 'making',
                    'sledge', 'facing'])

        result = list(exercise25(5, 6, m_trie, n_trie))
        self.assertEqual(len(result), 0)

        m_trie.add('stage')

        result = list(exercise25(5, 6, m_trie, n_trie))
        self.assertEqual(len(result), 1)
        self.assertEqual(result, [('slums', 'total', 'awoke', 'tepid', 'using',
                                  'stage')])

        result = list(exercise25(2, 2))
        self.assertEqual(len(result), 2177)

    def test_long_exercise25(self):
        # 5 x 6
        result = list(islice(exercise25(5, 6), 191))

        self.assertEqual(result[0],
                         ('aargh', 'blare', 'lapin', 'atilt', 'tense',
                          'edged'))

        self.assertEqual(result[190],
                         ('abaca', 'baths', 'bites', 'elude', 'sines',
                          'seers'))

        # result = list(exercise25(5, 6))
        # self.assertEqual(len(result), 625415)

        # for word in exercise25(5, 6):
        #     print(word)


if __name__ == '__main__':
    unittest.main(exit=False)
