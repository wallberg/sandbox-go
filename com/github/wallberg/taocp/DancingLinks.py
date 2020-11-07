# -*- coding: utf-8 -*-
import unittest
from itertools import islice
from copy import deepcopy

'''
Explore Dancing Links from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2.1 Dancing Links
'''


def exact_cover(items, options, secondary=[], stats=None, progress=None):
    '''Algorithm X. Exact cover via dancing links.

    Arguments:
    items     -- sequence of primary items
    options   -- sequence of options; every option must contain at least one
                 primary item

    Keyword Arguments:
    secondary -- sequence of secondary items
    stats     -- dictionary to accumulate runtime statistics
    progress  -- display progress report, every 'progress' number of level
                 entries
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
        # Iterate over the options
        options = []
        for p in x:
            option = []
            # Move back to first element in the option
            while top[p-1] > 0:
                p -= 1
            # Iterate over elements in the option
            q = p
            while top[q] > 0:
                option.append(name[top[q]])
                q += 1
            options.append(tuple(option))
        return tuple(options)

    def mrv():
        ''' Minimum Remaining Values heuristic. '''
        nonlocal rlink, llen

        theta = -1
        p = rlink[0]
        while p != 0:
            lambd = llen[p]
            if lambd < theta or theta == -1:
                theta = lambd
                i = p
                if theta == 0:
                    return i

            p = rlink[p]

        return i

    def dump():
        nonlocal name, rlink, dlink

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

    def show_progress():
        nonlocal stats, level, n, max_level, x, top, size, llen

        est = 0.0  # estimate of percentage done
        tcum = 1

        print(f'Current level {level} of max {max_level}')

        # Iterate over the options
        for p in x[0:level]:
            # Cyclically gather the items in the option, beginning at p
            option = []
            q = p
            while True:
                option.append(str(name[top[q]]))
                q += 1
                if top[q] <= 0:
                    q = ulink[q]
                if q == p:
                    break

            # Get position stats for this option
            i = top[p]
            q = dlink[i]
            k = 1
            while q != p and q != i:
                q = dlink[q]
                k += 1

            if q != i:
                kstat = f'{k} of {llen[i]}'
                tcum *= llen[i]
                est += (k - 1) / tcum
            else:
                kstat = "not in this list"

            print(f'  {" ".join(option)} ({kstat})')

        est += 1 / (2 * tcum)

        print(f"  solutions={stats['solutions']}, nodes={stats['nodes']}, est={est:4.4f}")
        print('---')

    # X1 [Initialize.]

    n1 = len(items)  # number of primary items
    n2 = len(secondary)  # number of secondary items
    n = n1 + n2  # total number of items

    if progress is not None:
        if stats is None:
            # Track progress the stats dictionary
            stats = {}
        assert isinstance(progress, int)
        theta, delta = progress, progress
        progress = True
        max_level = -1

    if stats is not None:
        stats['level_count'] = [0] * n
        stats['nodes'] = 0
        stats['solutions'] = 0

    # Fill out the item tables
    name = [None] * (n + 2)
    llink = [None] * (n + 2)
    rlink = [None] * (n + 2)

    i = 0
    for item in list(items) + list(secondary):
        i += 1
        name[i] = item
        llink[i] = i - 1
        rlink[i-1] = i

    # two doubly linked lists, primary and secondary
    # head of the primary list is at i=0
    # head of the secondary list is at i=n+1
    llink[n+1] = n
    rlink[n] = n + 1
    llink[n1+1] = n + 1
    rlink[n+1] = n1 + 1
    llink[0] = n1
    rlink[n1] = 0

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

    # dump()

    goto = 'X2'
    while True:

        if goto == 'X2':
            # [Enter level l.]
            if stats is not None:
                stats['level_count'][level] += 1
                stats['nodes'] += 1

            if rlink[0] == 0:
                # visit the solution
                if stats is not None:
                    stats['solutions'] += 1
                yield solution(x[0:level])
                goto = 'X8'
            else:
                goto = 'X3'

            if progress:
                if level > max_level:
                    max_level = level
                if stats['nodes'] >= theta:
                    show_progress()
                    theta += delta

        elif goto == 'X3':
            # [Choose i.]
            i = mrv()
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
            if stats is not None:
                stats['nodes'] += 1

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


def langford_pairs(n, **kwargs):
    '''Return solutions for Langford pairs of n values.'''

    items = [i for i in range(1, n+1)] + [f's{j-1}' for j in range(1, 2*n+1)]

    options = []
    for i in range(1, n+1):
        j = 1
        k = j + i + 1
        while k <= 2*n:
            # Exercise 15: Omit the reversals
            if i != n - (n % 2 == 0) or j <= n / 2:
                options.append((i, f's{j-1}', f's{k-1}'))
            j += 1
            k += 1

    for solution in exact_cover(items, options, **kwargs):
        x = [None] * (2 * n)
        for option in solution:
            x[int(option[1][1])] = option[0]
            x[int(option[2][1])] = option[0]

        yield tuple(x)


def n_queens(n, **kwargs):
    '''Return solutions for the n-queens problem.'''

    items = []
    sitems = []
    options = []

    for i in range(1, n+1):
        row = 'r' + str(i)
        items.append(row)
        for j in range(1, n+1):
            col = 'c' + str(j)
            if i == n:
                items.append(col)
            up_diag = 'a' + str(i + j)
            down_diag = 'b' + str(i - j)
            if up_diag not in sitems:
                sitems.append(up_diag)
            if down_diag not in sitems:
                sitems.append(down_diag)

            options.append((row, col, up_diag, down_diag))

    for solution in exact_cover(items, options, secondary=sitems, **kwargs):
        yield tuple(option[:2] for option in solution)


def sudoku(grid, **kwargs):
    '''Return solutions for SuDoku puzzles.'''

    def build_option(i, j, k, x):
        '''Build the (p, r, c, b) option. '''

        return(('p' + str(i) + str(j),
                'r' + str(i) + str(k),
                'c' + str(j) + str(k),
                'x' + str(x) + str(k)))

    grid = [[int(n) for n in row] for row in grid]

    # Get the known items provided in grid parameter
    known_items = set()
    for i in range(9):  # row number
        for j in range(9):  # column number
            k = grid[i][j]  # cell value
            if k > 0:
                x = 3 * (i // 3) + (j // 3)  # box number
                for item in build_option(i, j, k, x):
                    known_items.add(item)

    items = set()
    options = []

    for i in range(9):  # row number
        for j in range(9):  # column number
            x = 3 * (i // 3) + (j // 3)  # box number
            for k in range (1, 10):  # cell value
                option = build_option(i, j, k, x)
                if all(item not in known_items for item in option):
                    items.update(option)
                    options.append(option)

    for solution in exact_cover(list(items), options, **kwargs):
        sln = deepcopy(grid)
        for p, r, _, _ in solution:
            i, j, k = int(p[1]), int(p[2]), int(r[2])
            sln[i][j] = k

        yield [''.join(str(k) for k in row) for row in sln]


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
                         [(3, 1, 2, 1, 3, 2)])

        result = sum(1 for s in langford_pairs(7))
        self.assertEqual(result, 26)

        result = sum(1 for s in langford_pairs(8))
        self.assertEqual(result, 150)

        result = sum(1 for s in langford_pairs(10))
        self.assertEqual(result, 0)

    def test_long_langford_pairs(self):
        result = sum(1 for s in langford_pairs(11))
        self.assertEqual(result, 17792)

    def test_n_queens(self):
        result = list(n_queens(4))
        self.assertEqual(result,
                         [(('r1', 'c2'), ('r2', 'c4'),
                           ('r3', 'c1'), ('r4', 'c3')),
                          (('r1', 'c3'), ('r2', 'c1'),
                           ('r3', 'c4'), ('r4', 'c2'))])

    def test_sudoku(self):
        pbm1 = ["083921657",
                "967345821",
                "251876493",
                "548132970",
                "729564138",
                "136798245",
                "372689514",
                "814253769",
                "695417380"]

        sln1 = ["483921657",
                "967345821",
                "251876493",
                "548132976",
                "729564138",
                "136798245",
                "372689514",
                "814253769",
                "695417382"]

        pbm2 = ["300200000",
                "000107000",
                "706030500",
                "070009080",
                "900020004",
                "010800050",
                "009040301",
                "000702000",
                "000008006"]

        sln2 = ["351286497",
                "492157638",
                "786934512",
                "275469183",
                "938521764",
                "614873259",
                "829645371",
                "163792845",
                "547318926"]

        # 29a
        pbm3 = ["003010000",
                "415000090",
                "206500300",
                "500080009",
                "070900032",
                "038004060",
                "000260403",
                "000300008",
                "320007950"]

        sln3 = ['793412685',
                '415638297',
                '286579314',
                '562183749',
                '174956832',
                '938724561',
                '859261473',
                '647395128',
                '321847956']

        # 29b
        pbm4 = ["000000300",
                "100400000",
                "000000105",
                "900000000",
                "000002600",
                "000053000",
                "050800000",
                "000900070",
                "083000040"]

        sln4 = ['597218364',
                '132465897',
                '864379125',
                '915684732',
                '348792651',
                '276153489',
                '659847213',
                '421936578',
                '783521946']

        result = list(sudoku(pbm1))
        self.assertEqual(result, [sln1])

        result = list(sudoku(pbm2))
        self.assertEqual(result, [sln2])

        result = list(sudoku(pbm3))
        self.assertEqual(result, [sln3])

        result = list(sudoku(pbm4))
        self.assertEqual(result, [sln4])


if __name__ == '__main__':
    unittest.main(exit=False)
