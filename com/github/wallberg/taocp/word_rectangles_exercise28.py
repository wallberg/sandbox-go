# -*- coding: utf-8 -*-
import unittest
from itertools import islice

from com.github.wallberg.taocp.Trie import PrefixTrie, CompressedPrefixTrie
from com.github.wallberg.taocp.Backtrack import basic_backtrack

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming
'''


def word_rectangles(m, n, m_trie=None, n_trie=None):
    '''
    Find m x n word rectangles using Basic Backtrack, using a prefix trie
    as in Exercise 24 and compressed prefix tree as in Exercise 28.
    '''

    # Load m length prefix trie
    if m_trie is None:
        m_trie = CompressedPrefixTrie()
        if m == 5:
            m_trie.load_sgb_words()
        else:
            m_trie.load_ospd4_words(m)
    else:
        assert isinstance(m_trie, CompressedPrefixTrie)

    # Load n length prefix trie
    if n_trie is None:
        n_trie = PrefixTrie()
        if n == 5:
            n_trie.load_sgb_words()
        else:
            n_trie.load_ospd4_words(n)
    else:
        assert isinstance(n_trie, PrefixTrie)

    mn = m * n

    # m_a is an m x n lookup table for n_trie nodes corresponding to the
    # prefixes of the first l columns of partial solution x
    n_a = [[0] * m for i in range(n+1)]

    # B1 [Initialize]
    level = 1
    i = 0
    j = 0
    x = [None] * mn

    domains = [0] * mn

    goto = 'B2'
    while True:

        if goto == 'B2':
            # [Enter level l.]

            if level > mn:
                yield tuple(''.join([chr(c+96) for c in x[i:i+m]])
                            for i in range(0, m*n, m))
                goto = 'B5'
            else:
                # Set min D_l
                node = 0 if j == 0 else domains[level-2][1]
                first_link = m_trie.trie[node]
                domains[level-1] = first_link
                x[level-1] = first_link[0]

                goto = 'B3'

        elif goto == 'B3':
            # [Try x_l.]

            goto = 'B4'

            c = x[level-1]

            n_node = n_trie.trie[n_a[i][j]][c]
            if n_node != 0:
                n_a[i+1][j] = n_node

                # P_l holds true, we can advance
                level += 1
                j += 1
                if j == m:
                    i += 1
                    j = 0

                goto = 'B2'

        elif goto == 'B4':
            # [Try again.]

            # Set next value of D_l
            next_link = domains[level-1][2]
            if next_link[0] == 0:
                # max D_l
                goto = 'B5'
            else:
                # next D_l
                domains[level-1] = next_link
                x[level-1] = next_link[0]
                goto = 'B3'

        elif goto == 'B5':
            # [Backtrack.]

            level -= 1
            j -= 1
            if j == -1:
                i -= 1
                j = m - 1

            if level > 0:
                if level < mn:
                    x[level] = None
                goto = 'B4'
            else:
                return


class Test(unittest.TestCase):

    def test_word_rectangles(self):

        # 2 x 3 simple
        m_trie = CompressedPrefixTrie()
        m_trie.add(['ab', 'cd', 'ef'])
        m_trie.add(['ag', 'ah', 'ai', 'ej', 'ek'])

        n_trie = PrefixTrie()
        n_trie.add(['ace', 'bdf'])
        n_trie.add(['alm', 'acn', 'bop', 'qrs'])

        result = list(word_rectangles(2, 3, m_trie, n_trie))
        self.assertEqual(result, [('ab', 'cd', 'ef')])

        # 5 x 6 simple
        m_trie = CompressedPrefixTrie()
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

        # 5 x 6 full
        result = list(islice(word_rectangles(5, 6), 191))

        self.assertEqual(result[0],
                         ('aargh', 'blare', 'lapin', 'atilt', 'tense',
                          'edged'))

        self.assertEqual(result[190],
                         ('abaca', 'baths', 'bites', 'elude', 'sines',
                          'seers'))

        # result = list(word_rectangles(5, 6))
        # self.assertEqual(len(result), 625415)

        # for solution in word_rectangles(5, 6):
        #     print(solution)


if __name__ == '__main__':
    unittest.main(exit=False)
