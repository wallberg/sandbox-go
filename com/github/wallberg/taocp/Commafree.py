# -*- coding: utf-8 -*-
import unittest
from math import floor
from itertools import product
from array import array

from com.github.wallberg.taocp.Backtrack import walkers_backtrack
from com.github.wallberg.taocp.Tuples import preprimes

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


def exercise35(word, words, n=None):
    '''
    Test if word is commafree with respect to the accepted words of size m
    '''

    if n is None:
        n = len(word)

    test_words = set(words + [word])
    for test1, test2 in product(test_words, test_words):
        testseq = test1 + test2

        for start in [0, 1]:
            count = 0
            for i in range(0, n):
                if testseq[start+i:start+i+n] in test_words:
                    count += 1
            if count > 1:
                return False

    return True


def commafree_classes(m, n):
    '''
    Find the largest commafree subset of m-ary tuples of length n.
    Use cyclic classes to group options which must contribute only one of their
    values.
    '''

    def S(xn, level, x):
        ''' Return values at level, which hold true for x_1 to x_(level-1) '''

        if level == 1:
            return domain

        elif x[level-2] is None:
            return [None]

        # Record the class of x_(l-1)
        x_class[level-2] = class_map[x[level-2]]
        # print(f'{x_class[:level-1]=}')

        # Check all words after x_(l-1)
        result = []
        i = domain.index(x[level-2]) + 1
        while i < len(domain):

            word = domain[i]
            # print(f'checking {word} with clas {class_map[word]}')

            # Ensure this word's class isn't included already
            if class_map[word] not in x_class[:level-1]:

                # Check if word is commafree with respect to x_1 .. x_(l-1)
                if exercise35(word, x[:level-1], n):
                    result.append(word)

            i += 1

        if len(result) > 0:
            # print(f'{level=} , {x[:level-1]=}')
            return result
        else:
            return [None]

    # Maximum length of commafree subset
    xn = (m**4 - m**2) // 4

    # classes of word cycles
    classes = list(word for word, j in preprimes(m, n) if j == n)
    # print(f'{xn=}, {classes=}')

    domain = []
    class_map = {}
    for clas in classes:
        word = clas
        for _ in range(n):
            domain.append(word)
            class_map[word] = clas
            word = word[1:n] + word[0:1]

    # Track the word class at each x_j
    x_class = [None] * xn

    # Generate all commafree subsets
    max_subset = []
    for cf in walkers_backtrack(xn, S):
        try:
            size = cf.index(None)
        except ValueError:
            size = n

        if size > len(max_subset):
            max_subset = cf[:size]

            # Stop at first subset to reach the maximum possible size
            if len(max_subset) == xn:
                break

    return max_subset

def commafree_four(m, g, max=0):
    '''
    Algorithm C. Four-letter commafree codes.

    m - alphabet size
    g - goal number of words in the code
    '''

    RED = 0
    BLUE = 1
    GREEN = 2

    def alpha(word):
        ''' Return integer representation of the code word. '''
        result = 0
        for letter in word:
            result *= m
            result += letter

        return result

    def tostring():
        ''' String representation of MEM table. '''

        table = []
        for i in range(22):
            row = []
            for j in range(M4):
                if i == 0:
                    if MEM[j] == RED:
                        row.append(' RED')
                    elif MEM[j] == BLUE:
                        row.append('BLUE')
                    else:
                        row.append('GREN')
                else:
                    alf = MEM[i*M4+j]
                    if alf == 0:
                        row.append('    ')
                    elif i % 3 == 2:
                        row.append(''.join(str(c) for c in ALF[alf]))
                    else:
                        row.append('{:-4x}'.format(MEM[i*M4 + j]))

            table.append(' | '.join(row) + '\n')

            if (i+1) % 3 == 1:
                table.append('-' * (7*M4-3) + '\n')

        return ''.join(table)

    def initialize_lists():

        # Initialize colors to RED
        for alf in range(M4):
            MEM[alf] = RED

        # Iterate over word classes
        for cl, clas in enumerate(word for word, j in preprimes(m, 4)
                                  if j == 4):

            # Iterate through the 4-cycle of words in this class
            word = clas
            for _ in range(4):
                alf = alpha(word)
                ALF[alf] = word

                # Skip 0100 and 1000 since they will generate symmetric
                # duplicates
                if word != (0, 1, 0, 0) and word != (1, 0, 0, 0):

                    MEM[alf] = BLUE

                    # Insert into 3 prefix and 3 suffix lists
                    offset = P1OFF
                    for ps in [alpha(word[0:1] + (0, 0, 0)),
                            alpha(word[0:2] + (0, 0)),
                            alpha(word[0:3] + (0, )),
                            alpha(word[3:4] + (0, 0, 0)),
                            alpha(word[2:4] + (0, 0)),
                            alpha(word[1:4] + (0, )),
                            ]:
                        tail = offset+M4+ps

                        if MEM[tail] == 0:
                            MEM[tail] = offset+ps

                        insert(alf, tail, offset-M4)

                        offset += 3*M4

                    # Insert into CLOFF
                    tail = CLOFF + M4 + (4 * cl)

                    if MEM[tail] == 0:
                        MEM[tail] = CLOFF+4*cl

                    insert(alf, tail, CLOFF-M4)

                # Cycle
                word = word[1:4] + word[0:1]

        print(tostring())

    def insert(alf, tail, ihead):
        ''' Insert a value into the list and the inverted list. '''

        MEM[MEM[tail]] = alf
        MEM[ihead+alf] = MEM[tail]
        MEM[tail] += 1

    # C1. [Initialize.]
    assert 2 <= m <= 7

    M2 = m**2
    M4 = m**4
    L = (M4 - M2) // 4

    assert L - m * (m - 1) <= g <= L

    M = floor(23.5 * M4)
    P1OFF = 2 * M4
    P2OFF = 5 * M4
    P3OFF = 8 * M4
    S1OFF = 11 * M4
    S2OFF = 14 * M4
    S3OFF = 17 * M4
    CLOFF = 20 * M4

    STAMP = [0] * M
    X = [0] * (L+1)
    C = [0] * (L+1)
    S = [0] * (L+1)
    U = [0] * (L+1)
    FREE = [0] * L
    IFREE = [0] * L
    UNDO = []
    sigma = 0

    # alpha to code word lookup table
    ALF = [0] * (16*3 * M)

    # Main table of lists: alpha, P1, P2, P3, S1, S3, S3, CL, POISON
    MEM = array('I', [0] * M)

    POISON = 22 * M4
    PP = POISON - 1
    MEM[PP] = POISON

    level = 1
    x = 1  # trial word
    c = 0  # trial word's class
    s = L - g  # "slack"
    f = 0  # number of free classes
    u = 0  # size of the UNDO stack

    # Fill in the tables
    initialize_lists()


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

    def test_commafree_classes(self):
        result = commafree_classes(2, 4)
        self.assertEqual(result, ((0, 0, 0, 1), (0, 0, 1, 1), (0, 1, 1, 1)))

        result = commafree_classes(3, 3)
        self.assertEqual(result, ((0, 0, 1), (0, 0, 2), (1, 1, 0), (2, 0, 1),
                                  (2, 1, 0), (2, 0, 2), (1, 1, 2), (2, 1, 2)))

if __name__ == '__main__':
    # unittest.main(exit=False)

    commafree_four(2, 3)