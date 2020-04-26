# -*- coding: utf-8 -*-
import unittest
import os

'''
Explore Trie Digital Searching, from The Art of Computer Programming, Volume 3,
Sorting and Searching, Second Edition, 1998.

§6.3 Digital Searching
'''


class WordTrie():

    def __init__(self, words=None):
        '''
        Initialize a new trie.

        If words is not None, use it a sequence of initial values.
        '''
        self.trie = []
        self.count = 0

        if words is not None:
            self.add(words)

    def add(self, word):
        ''' Add a single word or sequence of words. '''

        if not isinstance(word, str):
            for w in word:
                self.add(w)
        else:
            self.add_single(word)

    def add_single(self, word):
        '''
        Add a word to the trie.  word will be downcased and should only
        contains letters a-z.  word can also be a sequence of words to add.
        '''

        if not isinstance(word, str):
            for w in word:
                self.add_single(w)
            return

        word = word.lower()

        k = 0
        li = 0
        while li < len(word):
            # Get the letter index in the node
            i = ord(word[li]) - 96

            # Add a node if necessary
            if k == len(self.trie):
                self.trie.append([0] * 27)

            # There are three possible states of self.trie[k][i]
            #  - Collision: a word is already there
            #  - Empty: put the word there
            #  - Pointer: follow the pointer to the next node

            if isinstance(self.trie[k][i], str):
                # Collision, there's a word already there
                word2 = self.trie[k][i]

                if word == word2:
                    # Word already in the trie, we're done
                    return

                # Expand nodes while word and word2 letters match
                while True:

                    # Add new node
                    next_k = len(self.trie)
                    self.trie.append([0] * 27)

                    self.trie[k][i] = next_k
                    k = next_k
                    li += 1

                    # Check if we have exhausted our common letters
                    if li == len(word) \
                       or li == len(word2) \
                       or word[li] != word2[li]:
                        break

                    i = ord(word[li]) - 96

                # Add word2
                if li == len(word2):
                    self.trie[k][0] = word2
                else:
                    i = ord(word2[li]) - 96
                    self.trie[k][i] = word2

                # Add word
                if li == len(word):
                    self.trie[k][0] = word
                else:
                    i = ord(word[li]) - 96
                    self.trie[k][i] = word

                self.count += 1
                return

            elif self.trie[k][i] == 0:
                # Empty, add the word here
                self.trie[k][i] = word
                self.count += 1
                return

            else:
                # Pointer, traverse to the next node
                k = self.trie[k][i]

            li += 1

        # We've exhausted the letters, so it goes in this node
        self.trie[k][0] = word

    def traverse(self):
        ''' Generator for all words in lexicographic order. '''

        def traverse_r(k):

            if k == len(self.trie):
                return

            for i in range(0, 27):

                if isinstance(self.trie[k][i], str):
                    yield(self.trie[k][i])

                elif self.trie[k][i] > 0:
                    for word in traverse_r(self.trie[k][i]):
                        yield word

        for word in traverse_r(0):
            yield word

    def __iter__(self):
        ''' Create an iterator over words in this trie. '''
        self.iterator = self.traverse()
        return self

    def __next__(self):
        ''' Next word from the iterator. '''
        return next(self.iterator)

    def __contains__(self, word):
        ''' Determine if word is in the trie. '''

        word = word.lower()

        k = 0
        li = 0
        while li < len(word):

            # Get the letter index in the node
            i = ord(word[li]) - 96

            if isinstance(self.trie[k][i], str):
                # Check for word match
                return word == self.trie[k][i]

            elif self.trie[k][i] == 0:
                # Not found
                return False

            else:
                # Advance to next node and letter
                k = self.trie[k][i]

            li += 1

        # We've exhausted the letters
        return self.trie[k][0] == word

    def __len__(self):
        ''' Number of words in the trie. '''
        return(self.count)

    def load_sgb_words(self):
        ''' Load Stanford GraphBase words into the trie. '''

        # Read in every line
        fn = os.path.join(os.path.dirname(__file__), 'data/sgb-words.txt')
        with open(fn) as f:
            for word in [line.strip() for line in f]:
                self.add_single(word)

    def load_ospd4_words(self, n=None):
        '''
        Load Official SCRABBLE Players Dictionary, 4th edition, words into the
        trie. Optionally limit to words of length n.
        '''

        # Read in every line
        fn = os.path.join(os.path.dirname(__file__), 'data/ospd4.txt')
        with open(fn) as f:
            for word in [line.strip() for line in f]:
                if n is None or len(word) == n:
                    self.add_single(word)


class PrefixTrie(WordTrie):
    '''
    A WordTrie which stores it's word in the full prefix path for all the
    letters of the word.  Since it doesn't stop at the first available node
    this implementation requires more memory.
    '''

    def add_single(self, word):
        '''
        Add a word to the trie.  word will be downcased and should only
        contains letters a-z.  word can also be a sequence of words to add.
        '''

        word = word.lower()

        k = 0
        li = 0
        word_existing = None

        while li < len(word):
            # Get the letter index in the node
            i = ord(word[li]) - 96

            # Add a node if necessary
            if k == len(self.trie):
                self.trie.append([0] * 27)

            # If there was a previous collision, now reposition that
            # existing word
            if word_existing:
                self.trie[k][0] = word_existing
                word_existing = None

            # There are three possible states of self.trie[k][i]
            #  - Collision: a word is already there
            #  - Empty: available for a new pointer or the word
            #  - Pointer: follow the pointer to the next node

            if isinstance(self.trie[k][i], str):
                # Collision, there's a word already there
                word_existing = self.trie[k][i]

                if word == word_existing:
                    # Word already in the trie, we're done
                    return

                # The existing word must be a shorter prefix to the
                # incoming word.

                # Add pointer to a new node
                next_k = len(self.trie)
                self.trie[k][i] = next_k
                k = next_k

            elif self.trie[k][i] == 0:
                # Empty slot
                if li == len(word) - 1:
                    # Last letter, put the word here
                    self.trie[k][i] = word
                else:
                    # Add pointer to a new node
                    next_k = len(self.trie)
                    self.trie[k][i] = next_k
                    k = next_k

            else:
                # Pointer, traverse to the next node
                k = self.trie[k][i]

                if li == len(word) - 1:
                    # Last letter, put the word here, in the 0 slot
                    self.trie[k][0] = word

            li += 1

        self.count += 1


class CompressedPrefixTrie(PrefixTrie):
    '''
    A Compressed PrefixTrie which uses an array of linked lists for storage,
    instead of a two dimensional array.  This makes letter traversal faster
    and saves space.  It also assumes that all words are of the same length, so
    there will be no prefix collisions.

    Link are in alpha order, with format: [letter, node, link]

    - letter is 1-26 ('a' -'z'), or 0 if last link in the list
    - node is node number of next letter, or 0 if last letter in the word
    - link is pointer to next letter at this node, or 0 if final link
    '''

    def add_single(self, word):
        '''
        Add a word to the trie.  word will be downcased and should only
        contains letters a-z.  word can also be a sequence of words to add.
        '''

        word = word.lower()

        node = 0

        for li in range(len(word)):
            # Get the letter
            c = ord(word[li]) - 96

            # Add a node if necessary
            if node == len(self.trie):
                self.trie.append([0, 0, 0])

            # Search the linked list to either find the existing entry or
            # insert a new one
            link = self.trie[node]
            while True:

                if link[0] == 0 or link[0] > c:
                    # Insert here
                    new_link = link.copy()
                    link[0] = c
                    link[2] = new_link

                    # Create new node, if necessary
                    if li == len(word) - 1:
                        link[1] = 0
                    else:
                        node = len(self.trie)
                        link[1] = node

                    break

                elif link[0] == c:
                    # Link already exists
                    node = link[1]

                    break

                # Advance to next link
                link = link[2]

        self.count += 1

    def traverse(self):
        ''' Generator for all words in lexicographic order. '''

        def traverse_r(node):

            if node == len(self.trie):
                return

            link = self.trie[node]
            while link[0] != 0:
                letter = chr(link[0] + 96)

                if link[1] == 0:
                    yield letter
                else:
                    for x in traverse_r(link[1]):
                        yield letter + x

                # Advance to next link
                link = link[2]

        for word in traverse_r(0):
            yield word

    def __contains__(self, word):
        ''' Determine if word is in the trie. '''

        raise NotImplemented()


class WordTrieTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.sgb = WordTrie()
        cls.sgb.load_sgb_words()

        cls.ospd4 = WordTrie()
        cls.ospd4.load_ospd4_words()

    def test_add_and_traverse(self):
        trie = WordTrie(['not', 'you', 'a', 'and', 'are', 'as', 'at', 'be',
                         'but', 'by', 'for', 'from', 'his', 'in', 'is', 'it',
                         'of', 'on', 'or', 'to', 'was', 'which', 'with', 'had',
                         'have', 'he'])

        trie.add('her')
        trie.add(['that', 'the', 'this'])

        result = list(word for word in trie)
        self.assertEqual(result, ['a', 'and', 'are', 'as', 'at', 'be', 'but',
                                  'by', 'for', 'from', 'had', 'have', 'he',
                                  'her', 'his', 'in', 'is', 'it', 'not', 'of',
                                  'on', 'or', 'that', 'the', 'this', 'to',
                                  'was', 'which', 'with', 'you'])

    def test_contains(self):
        trie = WordTrie()
        trie.add('a')
        trie.add('aa')
        trie.add('aaa')
        trie.add('b')

        self.assertTrue('a' in trie)
        self.assertTrue('aa' in trie)
        self.assertTrue('aaa' in trie)
        self.assertTrue('b' in trie)

        self.assertTrue('aaaa' not in trie)
        self.assertTrue('c' not in trie)
        self.assertTrue('ab' not in trie)

        self.assertTrue('staph' in self.sgb)
        self.assertTrue('skewz' not in self.sgb)

    def test_len(self):
        trie = WordTrie()
        self.assertEqual(len(trie), 0)

        trie.add('a')
        self.assertEqual(len(trie), 1)

        trie.add('b')
        trie.add('a')
        self.assertEqual(len(trie), 2)

        self.assertEqual(len(self.sgb), 5757)

        self.assertEqual(len(self.ospd4), 178379)

    def test_load_sgb_words(self):
        words = list([word for word in self.sgb])
        self.assertEqual(len(words), 5757)

        self.assertEqual(words[0], 'aargh')
        self.assertEqual(words[1], 'abaca')
        self.assertEqual(words[2], 'abaci')

        self.assertEqual(words[428], 'berry')
        self.assertEqual(words[1248], 'deque')
        self.assertEqual(words[2968], 'mails')
        self.assertEqual(words[4458], 'skews')
        self.assertEqual(words[4733], 'staph')

        self.assertEqual(words[5754], 'zooks')
        self.assertEqual(words[5755], 'zooms')
        self.assertEqual(words[5756], 'zowie')

    def test_load_ospd4_words(self):
        trie = WordTrie()

        trie.load_ospd4_words(2)
        self.assertEqual(len(trie), 101)

        trie.load_ospd4_words(15)
        self.assertEqual(len(trie), 3258)


class PrefixTrieTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.sgb = PrefixTrie()
        cls.sgb.load_sgb_words()

        cls.ospd4 = PrefixTrie()
        cls.ospd4.load_ospd4_words()

    def test_add(self):
        ptrie = PrefixTrie(['a', 'aa', 'aaaaa', 'abbb', 'abb', 'ab', 'a'])

        result = list(word for word in ptrie)
        self.assertEqual(result, ['a', 'aa', 'aaaaa', 'ab', 'abb', 'abbb'])

        ptrie = PrefixTrie(['not', 'you', 'a', 'and', 'are', 'as', 'at', 'be',
                            'but', 'by', 'for', 'from', 'his', 'in', 'is',
                            'it', 'of', 'on', 'or', 'to', 'was', 'which',
                            'with', 'had', 'have', 'he'])

        ptrie.add('her')
        ptrie.add(['that', 'the', 'this'])

        result = list(word for word in ptrie)
        self.assertEqual(result, ['a', 'and', 'are', 'as', 'at', 'be', 'but',
                                  'by', 'for', 'from', 'had', 'have', 'he',
                                  'her', 'his', 'in', 'is', 'it', 'not', 'of',
                                  'on', 'or', 'that', 'the', 'this', 'to',
                                  'was', 'which', 'with', 'you'])

    def test_load_sgb_words(self):
        words = list([word for word in self.sgb])
        self.assertEqual(len(words), 5757)

        self.assertEqual(words[0], 'aargh')
        self.assertEqual(words[1], 'abaca')
        self.assertEqual(words[2], 'abaci')

        self.assertEqual(words[428], 'berry')
        self.assertEqual(words[1248], 'deque')
        self.assertEqual(words[2968], 'mails')
        self.assertEqual(words[4458], 'skews')
        self.assertEqual(words[4733], 'staph')

        self.assertEqual(words[5754], 'zooks')
        self.assertEqual(words[5755], 'zooms')
        self.assertEqual(words[5756], 'zowie')

    def test_load_ospd4_words(self):
        trie = WordTrie()

        trie.load_ospd4_words(2)
        self.assertEqual(len(trie), 101)

        trie.load_ospd4_words(15)
        self.assertEqual(len(trie), 3258)


class CompressedPrefixTrieTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.sgb = CompressedPrefixTrie()
        cls.sgb.load_sgb_words()

        cls.ospd4 = CompressedPrefixTrie()
        cls.ospd4.load_ospd4_words(6)

    def test_load_sgb_words(self):
        self.assertEqual(len(self.sgb), 5757)

        # for node in self.sgb.trie:
        #     print(node)

        words = list([word for word in self.sgb])
        self.assertEqual(len(words), 5757)

        self.assertEqual(words[0], 'aargh')
        self.assertEqual(words[1], 'abaca')
        self.assertEqual(words[2], 'abaci')

        self.assertEqual(words[428], 'berry')
        self.assertEqual(words[1248], 'deque')
        self.assertEqual(words[2968], 'mails')
        self.assertEqual(words[4458], 'skews')
        self.assertEqual(words[4733], 'staph')

        self.assertEqual(words[5754], 'zooks')
        self.assertEqual(words[5755], 'zooms')
        self.assertEqual(words[5756], 'zowie')

    def test_load_ospd4_words(self):
        self.assertEqual(len(self.ospd4), 15727)


if __name__ == '__main__':
    unittest.main(exit=False)
