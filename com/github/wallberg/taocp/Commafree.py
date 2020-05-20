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

        nonlocal RED, BLUE, GREEN, M4, ALF, MEM, UNDO

        lists = {2:'P1', 5:'P2', 8:'P3', 11:'S1', 14:'S2', 17:'S3', 20:'CL'}

        table = []

        # Add header
        table.append('      ' + '   '.join([f'{j:-4x}' for j in range(M4)]))

        tail = M

        for i in range(24):
            row = [f'{i*16:-3x}']

            for j in range(M4):
                n = i*M4 + j

                if n < len(MEM):
                    if i == 0:
                        if MEM[j] == RED:
                            row.append(' RED')
                        elif MEM[j] == BLUE:
                            row.append('BLUE')
                        else:
                            row.append(' GRN')
                    else:
                        alf = MEM[n]
                        if alf == 0:
                            row.append('    ')
                        elif i % 3 == 2:
                            if MEM[n + M4] > 0:
                                tail = MEM[n + M4]
                            if n < tail:
                                row.append(''.join(str(c) for c in ALF[alf]))
                            else:
                                row.append('xxxx')
                        else:
                            row.append('{:-4x}'.format(MEM[n]))
            row.append(' ')

            if i in lists:
                row[-1] += lists[i]

            table.append(' | '.join(row))

            if (i+1) % 3 == 1:
                table.append('    |-' + '-' * (7*M4-2) + '|')

        table.append('')

        # table.append('UNDO=' + ', '.join(['{:x}:{:x}'.format(a, v) for a, v in UNDO[:u]]))

        return '\n'.join(table)

    def prefixes_suffixes(word):
        '''
        Return 6-tuple of prefixes and suffixes of word.
        Return in table order or pair order.
        '''

        p1 = alpha(word[0:1] + (0, 0, 0))
        p2 = alpha(word[0:2] + (0, 0))
        p3 = alpha(word[0:3] + (0, ))
        s1 = alpha(word[3:4] + (0, 0, 0))
        s2 = alpha(word[2:4] + (0, 0))
        s3 = alpha(word[1:4] + (0, ))

        return (p1, p2, p3, s1, s2, s3)

    def initialize_mem():

        nonlocal RED, BLUE, GREEN, M4, ALF, MEM, P1OFF, CLOFF

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
                    for ps in prefixes_suffixes(word):
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

    def store(a, v):
        ''' Store v at MEM[a]; save original value on the UNDO stack '''

        nonlocal MEM, STAMP, UNDO, u, sigma

        # print(f'store: {a=}, {v=}, {STAMP[a]=}, {sigma=}')

        # Check if MEM[a] has been changed yet this round
        if STAMP[a] != sigma:
            # No, indicate now that is has been
            STAMP[a] = sigma
            # Save a and original value on the UNDO stack
            UNDO[u] = (a, MEM[a])
            u += 1

        MEM[a] = v

    def get_word_from_class():
        '''
        Get a word from a class with the least number of blue words.
        Exercise 44.
        '''

        nonlocal M4, MEM, CLOFF, FREE, PP, POISON, x, c

        print('get_word_from_class')

        r = m + 1  # number of words in a class with least blue words

        # Iterate over free classes
        for k in range(f):
            t = FREE[k]  # a free class
            j = MEM[CLOFF + 4*t + M4] - (CLOFF + 4*t)  # Size of class list

            # Does this class have the fewest words seen so far?
            if j < r:
                r, cl = j, t
                if r == 0:
                    x = -1
                    break

        print(f'{r=}')

        if r > 0:
            # Set x to a word in the class
            x = MEM[CLOFF + 4*cl]

        if r > 1:
            # Use the poison list to find an x that maximizes the number of
            # blue words that could be killed on the other side of the prefix
            # or suffix that contains x.
            q = 0
            pp = MEM[PP]
            p = POISON
            while p < pp:
                y = MEM[p]  # head of poison prefix list
                z = MEM[p + 1]  # head of poison suffix list

                yp = MEM[y + M4]  # tail of poison prefix list
                zp = MEM[z + M4] # tail of poison suffix list

                if y == yp or z == zp:
                    # Delete entry p from the poison list
                    # p, pp, q = delete(p, pp, y, z, yp, zp, q)
                    pp -= 2
                    if p != pp:
                        store(p, MEM[pp])
                        store(p+1, MEM[pp+1])
                    else:
                        p += 2
                        ylen = yp - y
                        zlen = zp - z
                        if  ylen >= zlen and ylen > q:
                            q = ylen
                            x = MEM[z]
                        if  ylen < zlen and zlen > q:
                            q = zlen
                            x = MEM[y]

            store(PP, pp)
            c = cl

    def rem(alpha, delta, omicron):
        ''' Remove an item from a list. '''

        nonlocal M4, MEM

        p = delta + omicron  # head pointer ?
        q = MEM[p + M4] - 1  # tail pointer ?

        if q >= p:
            # list p isn't closed or being killed
            store(p + M4, q)
            t = MEM[alpha + omicron - M4]

            if t != q:
                y = MEM[q]
                store(t, y)
                store(y + omicron - M4, t)

    def close(delta, omicron):
        ''' Close list delta+omicron. '''

        nonlocal M4, MEM

        print(f'close: {delta=:x}, {omicron=:x}')

        p = delta + omicron  # head of the class list
        q = MEM[p + M4]  # tail of the class list

        print(f'{p=:x}, {q=:x}')

        # Check if already closed
        if q != p - 1:
            # Close by setting tail to head-1
            store(p + M4, p - 1)

        # Return the head and tail
        return (p, q)

    def red(alf, c):
        ''' Make alpha RED. '''

        nonlocal RED, P1OFF, M4, ALF

        print(f'red: {alf=}, {c=}')

        store(alf, RED)

        offset = P1OFF
        for ps in prefixes_suffixes(ALF[alf]):
            rem(alf, ps, offset)  # remove from pre- or suffix list
            offset += 3*M4
        rem(x, 4*c, offset)  # remove from class list

    def green(alf, c):
        ''' Make alpha GREEN. '''

        nonlocal GREEN, P1OFF, CLOFF, M4, ALF

        print(f'green: {alf=}, {c=}')

        store(alf, GREEN)
        print(tostring())

        offset = P1OFF
        for ps in prefixes_suffixes(ALF[alf]):
            close(ps, offset)  # close pre- or suffix list
            offset += 3*M4
        p, q = close(4*c, CLOFF)  # close class list

        # Close the other words in this class
        print(f'{p=}, {q=}')
        for r in range(p, q):
            if MEM[r] != x:
                red(MEM[r], c)
                print(tostring())

    # C1. [Initialize.]
    print("C1.")

    assert 2 <= m <= 7

    M2 = m**2
    M4 = m**4
    L = (M4 - M2) // 4  # number of word classes
    M = floor(23.5 * M4)  # size of the main table, MEM

    assert L - m * (m - 1) <= g <= L

    MEM = array('I', [0] * M) # color, P1, P2, P3, S1, S2, S3, CL, POISON

    P1OFF = 2 * M4  # offsets into MEM for P1, P2, P3, S1, S2, S3, CL
    P2OFF = 5 * M4
    P3OFF = 8 * M4
    S1OFF = 11 * M4
    S2OFF = 14 * M4
    S3OFF = 17 * M4
    CLOFF = 20 * M4

    POISON = 22 * M4  # head of the poison list
    PP = POISON - 1  # tail of the poison list
    MEM[PP] = POISON

    level = 1  # backtrack level and index into X, C, S, U

    u = 0  # size of the UNDO stack
    U = [0] * (L+1)
    UNDO = [None] * 10000  # UNDO stack

    STAMP = [-1] * M  # store MEM[a] in UNDO only once per sigma
    sigma = 0

    x = 1  # trial word
    X = [0] * (L+1)

    c = 0  # trial word's class, simple index into the class l
    C = [0] * (L+1)

    s = L - g  # "slack"
    S = [0] * (L+1)

    f = L  # number of free classes, aka tail pointer for the free class list
    FREE = [c for c in range(L)]
    IFREE = [c for c in range(L)]

    # alpha to code word lookup table
    ALF = [0] * (16*3 * M)

    # Fill in the main tables
    initialize_mem()

    # Begin the main event loop
    step = 'C2'
    while True:

        if step == 'C2':
            # [Enter level.]
            print(f'C2. {level=}')

            if level == L:
                yield tuple(X[0:level])
                step = 'C6'

            else:
                # Choose a candidate word x and class c
                get_word_from_class()
                print(tostring())
                print(f'{x=}')

                step = 'C3'

        elif step == 'C3':
            # [Try the candidate.]

            U[level] = u
            sigma += 1

            step = 'C4'
            if x < 0:
                if s == 0 or level == 0:
                    step = 'C6'
                else:
                    s -= 1
            else:
                # Make x green
                green(x, c)

                # Add the three prefix, suffix pairs to the poison list
                pp = MEM[PP] + 6

                p1, p2, p3, s1, s2, s3 = prefixes_suffixes(ALF[x])
                store(pp - 6, P1OFF + p1)
                store(pp - 5, S3OFF + s3)
                store(pp - 4, P2OFF + p2)
                store(pp - 3, S2OFF + s2)
                store(pp - 2, P3OFF + p3)
                store(pp - 1, S1OFF + s1)

                # # ??
                # p = POISON

                # # Iterate over poison prefix/suffix pairs
                # while p < pp:
                #     # MEM[p:p+2] is one poison prefix/suffix pair

                #     y = MEM[p]  # head of the prefix list
                #     z = MEM[p + 1]  # head of the suffix list

                #     yp = MEM[y + M4]  # tail of the prefix list
                #     zp = MEM[z + M4] # tail of the suffix list

                #     if y == yp or z == zp:
                #         # Delete entry p from the poison list
                #         pp -= 2
                #         if p != pp:
                #             store(p, MEM[pp])
                #             store(p+1, MEM[pp+1])
                #         else:
                #             p += 2
                #             ylen = yp - y
                #             zlen = zp - z
                #             if  ylen >= zlen and ylen > q:
                #                 q = ylen
                #                 x = MEM[z]
                #             if  ylen < zlen and zlen > q:
                #                 q = zlen
                #                 x = MEM[y]

                #     elif yp < y and zp < z:
                #         # A poisoned pair is present
                #         step = 'C6'

                #     elif yp > y and zp > z:
                #         p += 2

                #     elif yp < y and zp > z:
                #         store(z + M4, z)
                #         for r in range(z, zp):
                #             red(MEM[r], c)  # class?
                #             # delete poison entry p

                #     else:  # yp > y and zp < z
                #         store(y + M4, y)
                #         for r in range(y, yp):
                #             red(MEM[r], c)  # class?
                #             # delete poison entry p

                # store(PP, pp)

                print(tostring())

                return

        elif step == 'C4':
            # [Make the move.]

            X[level] = x
            C[level] = c
            S[level] = s

            # Delete class c from the active list
            p = IFREE[c]
            f -= 1

            if p != f:
                y = FREE[f]
                FREE[p] = y
                IFREE[y] = p
                FREE[f] = c
                IFREE[c] = f

            level += 1
            step = 'C2'

        elif step == 'C5':
            # [Try again.]

            while u > U[level]:
                u -= 1
                a, v = UNDO[u]
                MEM[a] = v

            sigma += 1

            # make x red
            red(x, c)

            step = 'C2'

        elif step == 'C6':
            # [Backtrack.]

            level -= 1

            if level == -1:
                return

            x = X[level]
            c = C[level]
            f += 1

            if x < 0:
                step = 'C6'  # repeat this step
            else:
                s = S[level]
                step = 'C5'

        else:
            raise Exception(f'Invalid Step: {step}')


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

    list(commafree_four(2, 3))
    # commafree_four(3, 18)
    # commafree_four(7, 588)