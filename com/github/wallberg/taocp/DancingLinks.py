# -*- coding: utf-8 -*-
import unittest
from itertools import islice

from com.github.wallberg.taocp.Trie import WordTrie, PrefixTrie

'''
Explore Dancing Links from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2.1 Dancing Links
'''


def exact_cover(items, options, stats=None):
    '''
    Algorithm X. Exact cover via dancing links.
    '''

    def hide(p):
        nonlocal top, ulink, dlink, llen

        q = p + 1
        while q != p:
            x = top[q]
            u, d = ulink[q], dlink[q]
            if x <= 0:
                q = u  # q was a spacer
            else:
                dlink[u], ulink[d] = d, u
                llen[x] -= 1
                q += 1

    def cover(i):
        nonlocal dlink, llink, rlink

        p = dlink[i]
        while p != i:
            hide(p)
            p = dlink[p]
        l, r = llink[i], rlink[i]
        rlink[l], llink[r] = r, l

    def unhide(p):
        nonlocal top, ulink, dlink, llen

        q = p - 1
        while q != p:
            x = top[q]
            u, d = ulink[q], dlink[q]
            if x <= 0:
                q = d  # q was a spacer
            else:
                dlink[u], ulink[d] = q, q
                llen[x] += 1
                q -= 1

    def uncover(i):
        nonlocal ulink, llink, rlink

        l, r = llink[i], rlink[i]
        rlink[l], llink[r] = i, i
        p = ulink[i]
        while p != i:
            unhide(p)
            p = ulink[p]

    def solution(x):
        options = []
        for p in x:
            option = []
            q = p
            while top[q] > 0:
                option.append(name[top[q]])
                q += 1
            options.append(tuple(option))
        return tuple(options)

    def dump():
        nonlocal name, rlink, dlink

        # print(n)
        # print(name)
        # print(llink)
        # print(rlink)
        # print(top)
        # print(ulink)
        # print(dlink)

        i = 0
        while rlink[i] != 0:
            i = rlink[i]
            item = [name[i]]
            x = i
            while dlink[x] != i:
                x = dlink[x]
                item.append(x)
            print(item)

        print('---')

    # X1 [Initialize.]

    n = len(items)

    # Fill out the item tables
    name = [None] * (n + 1)
    llink = [None] * (n + 1)
    rlink = [None] * (n + 1)

    i = 0
    for item in items:
        i += 1
        name[i] = item
        llink[i] = i - 1
        rlink[i-1] = i

    llink[0] = n
    rlink[n] = 0

    # Fill out the option tables
    n_options = len(options)
    n_optionitems = sum(len(option) for option in options)
    size = n + 1 + n_options + 1 + n_optionitems

    top = [None] * size
    llen = top  # first n+1 elements of top
    ulink = [None] * size
    dlink = [None] * size

    # Set empty list for each item
    for x in range(1, n+1):
        llen[x] = 0
        ulink[x] = x
        dlink[x] = x

    # Insert each of the options and their items
    x = n + 1
    spacer = 0
    top[x] = spacer
    spacer_x = x

    # Iterate over each option
    for option in options:
        # Iterate over each item in this option
        for item in option:
            x += 1
            i = name.index(item)
            top[x] = i

            # Insert into the option list for this item
            llen[i] += 1  # increase the size by one
            head = i
            tail = i
            while dlink[tail] != head:
                tail = dlink[tail]

            dlink[tail] = x
            ulink[x] = tail

            ulink[head] = x
            dlink[x] = head

        # Insert spacer at end of each option
        dlink[spacer_x] = x
        x += 1
        ulink[x] = spacer_x + 1

        spacer -= 1
        top[x] = spacer
        spacer_x = x

    z = size - 1
    level = 0
    x = [None] * n_options

    goto = 'X2'
    while True:

        if goto == 'X2':
            # [Enter level l.]
            if rlink[0] == 0:
                # visit the solution
                yield solution(x[0:level])
                goto = 'X8'
            else:
                goto = 'X3'

        elif goto == 'X3':
            # [Choose i.]
            # TODO: Use llen(i) MRV instead
            i = rlink[0]
            goto = 'X4'

        elif goto == 'X4':
            # [Cover i.]
            cover(i)
            x[level] = dlink[i]
            goto = 'X5'

        elif goto == 'X5':
            # [Try x_l.]
            if x[level] == i:
                goto = 'X7'
            else:
                p = x[level] + 1
                while p != x[level]:
                    j = top[p]
                    if j <= 0:
                        p = ulink[p]
                    else:
                        cover(j)
                        p += 1
                level += 1
                goto = 'X2'

        elif goto == 'X6':
            # [Try again.]
            p = x[level] - 1
            while p != x[level]:
                j = top[p]
                if j <= 0:
                    p = dlink[p]
                else:
                    uncover(j)
                    p -= 1
            i = top[x[level]]
            x[level] = dlink[x[level]]
            goto = 'X5'

        elif goto == 'X7':
            # [Backtrack.]
            uncover(i)
            goto = 'X8'

        elif goto == 'X8':
            # [Leave level l.]
            if level == 0:
                return
            else:
                level -= 1
                goto = 'X6'


def langford_pairs(n):
    ''' Return solutions for Langford pairs of n values. '''

    items = [i for i in range(1, n+1)] + [('s', j-1) for j in range(1, 2*n+1)]

    options = []
    for i in range(1, n+1):
        j = 1
        k = j + i + 1
        while k <= 2*n:
            options.append((i, ('s', j-1), ('s', k-1)))
            j += 1
            k += 1

    for solution in exact_cover(items, options):
        x = [None] * (2 * n)
        for option in solution:
            x[option[1][1]] = option[0]
            x[option[2][1]] = option[0]

        yield tuple(x)


EXAMPLE_6 = (('c', 'e'),
             ('a', 'd', 'g'),
             ('b', 'c', 'f'),
             ('a', 'd', 'f'),
             ('b', 'g'),
             ('d', 'e', 'g'))


class Test(unittest.TestCase):

    def test_exact_cover(self):
        result = list(exact_cover(('a', 'b', 'c', 'd', 'e', 'f', 'g'),
                                  EXAMPLE_6))
        self.assertEqual(result, [(('a', 'd', 'f'), ('b', 'g'), ('c', 'e'))])

    def test_langford_pairs(self):
        result = list(langford_pairs(3))
        self.assertEqual(result,
                         [(3, 1, 2, 1, 3, 2),
                          (2, 3, 1, 2, 1, 3)])

        result = sum(1 for s in langford_pairs(7))
        self.assertEqual(result, 52)

        result = sum(1 for s in langford_pairs(8))
        self.assertEqual(result, 300)

        result = sum(1 for s in langford_pairs(9))
        self.assertEqual(result, 0)

if __name__ == '__main__':
    unittest.main(exit=False)
