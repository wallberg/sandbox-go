# -*- coding: utf-8 -*-
import unittest
import logging

from com.github.wallberg.taocp.Backtrack import basic_backtrack

'''
Explore Backtrack Programming from The Art of Computer Programming, Volume 4,
Fascicle 5, Mathematical Preliminaries Redux; Introduction to Backtracking;
Dancing Links, 2020

§7.2.2 Backtrack Programming
'''

logger = logging.getLogger('com.github.wallberg.taocp.backtrack_exercise71')


def exercise_71_warmup(x, stats=None):
    '''
    Generate solutions to Exercise 71 warmup (Two Questions) with n correct
    answers.
    '''

    assert len(x) == 2

    def domain(n, k):
        return ('A', 'B')

    def property(n, level, a):
        '''
        Check for count of correct answers, in this order of questions:
         i | Q
        ---+---
         0 |  1
         1 |  2
        '''

        nonlocal x

        valid = True

        # a[0] - Question 1
        if level > 0:

            if a[0] == 'A' and level > 1:
                valid = (x[0] == (a[1] == 'B'))

            elif a[0] == 'B':
                valid = (x[0] == (a[0] == 'A'))


        # a[1] - Question 2
        if valid and level > 1:

            if a[1] == 'A':
                valid = (x[1] == x[0])

            elif a[1] == 'B':
                valid = (x[1] == ((not x[1]) != (a[0] == 'A')))

        return valid


    for a in basic_backtrack(2, domain, property, stats=stats):
        yield a


def exercise_71(grading, stats=None):
    '''
    Generate solutions to Exercise 71 (Twenty Questions) which match the
    True/False grading.
    '''

    assert len(grading) == 20

    DEBUG = logger.isEnabledFor(logging.DEBUG)
    INFO = logger.isEnabledFor(logging.INFO)

    def domain(n, k):
        return ('A', 'B', 'C', 'D', 'E')

    def property(n, level, a):
        '''
        Check for correctness of answers, returning False if we've
        reached a cutoff point.
        '''

        nonlocal x, qi, iq, primes, max_primes, INFO

        valid = True

        # Answer counts are used multiple times
        count = {'A': 0, 'B': 0, 'C': 0, 'D': 0, 'E': 0}
        n_false = 0
        for i in range(level):
            count[a[i]] += 1
            if x[i] is False:
                n_false += 1

        values = sorted(count.values())

        # target = ['D', 'C', 'C', 'A', 'E', 'A', 'D', 'A', 'A', 'C']
        # DEBUG = (a[:2] == target[:2])

        # a[0] - Question 3
        if level > 0:
            if DEBUG: print('q3')

            # Iterate by i over known answers
            i = 0
            while valid and i < level:
                q = iq[i]

                if q != 20:
                    j = qi[q+1]

                    # Check if we have an answer for q+1
                    if level > j:
                        if a[0] == 'A':
                            valid = (x[0] == ((a[i] == a[j]) == (q == 15)))
                        elif a[0] == 'B':
                            valid = (x[0] == ((a[i] == a[j]) == (q == 16)))
                        elif a[0] == 'C':
                            valid = (x[0] == ((a[i] == a[j]) == (q == 17)))
                        elif a[0] == 'D':
                            valid = (x[0] == ((a[i] == a[j]) == (q == 18)))
                        elif a[0] == 'E':
                            valid = (x[0] == ((a[i] == a[j]) == (q == 19)))

                i += 1

        # a[1] - Question 15
        i = qi[15]
        if valid and level > i:
            if DEBUG: print('q15')

            if a[i] in ('A', 'E'):
                valid = (x[i] == False)
            else:
                for j in range(level):
                    if iq[j] % 2 == 1:

                        if a[i] == 'B':
                            valid = (x[i] == ((a[j] == 'A') == (iq[j] == 9)))
                        elif a[i] == 'C':
                            valid = (x[i] == ((a[j] != 'A') or (iq[j] != 11)))
                        elif a[i] == 'D':
                            valid = (x[i] == ((a[j] == 'A') == (iq[j] == 13)))

                    if not valid:
                        break


        # a[2] - Question 20
        if valid and level > 2:
            if DEBUG: print('q20')

            if a[2] in ('D', 'E'):
                valid = (x[2] == False)
            elif a[2] == 'A':
                valid = (x[2] == (n_false <= 2))
            elif a[2] == 'B':
                valid = (x[2] == (n_false <= 1))
            elif a[2] == 'C':
                valid = (x[2] == (n_false == 0))

        # a[3] - Question 19
        if valid and level > 3:
            if DEBUG: print('q19')

            if a[3] == 'B':
                valid = (x[3] == False)
            else:
                # Iterate over questions 14-20
                last_B = 0
                for q in range(14, 21):
                    i = qi[q]
                    if level > i and a[i] == 'B':
                        last_B = q

                if a[3] == 'A':
                    valid = (x[3] == ((level <= qi[14]) or (last_B == 14)))
                elif a[3] == 'C':
                    valid = (x[3] == ((level <= qi[16]) or (last_B == 16)))
                elif a[3] == 'D':
                    valid = (x[3] == ((level <= qi[17]) or (last_B == 17)))
                elif a[3] == 'E':
                    valid = (x[3] == ((level <= qi[18]) or (last_B == 18)))

        # a[4] - Question 2
        i = qi[2]
        if valid and level > i:
            # if DEBUG: print('q2')

            # Iterate over questions 3-12
            next_same = 0
            for q in range(3, 13):
                j = qi[q]
                # if DEBUG: print(f'{level=}, {i=}, {a[i]=}, {q=}, {j=}, {a[j]=}')
                if level > j and a[i] == a[j]:
                    next_same = q
                    break

            # if DEBUG: print(f'{next_same=}')

            if a[i] == 'A':
                valid = (x[i] == ((level <= qi[4]) or (next_same == 4)))
            elif a[i] == 'B':
                valid = (x[i] == ((level <= qi[6]) or (next_same == 6)))
            elif a[i] == 'C':
                valid = (x[i] == ((level <= qi[8]) or (next_same == 8)))
            elif a[i] == 'D':
                valid = (x[i] == ((level <= qi[10]) or (next_same == 10)))
            elif a[i] == 'E':
                valid = (x[i] == ((level <= qi[12]) or (next_same == 12)))

        # a[5] - Question 1
        if valid and level > 5:
            if DEBUG: print('q1')

            if a[5] == 'A':
                valid = (x[5] == True)
            else:
                # Iterate over questions 2-5
                first_A = 0
                for q in range(2, 6):
                    i = qi[q]
                    if level > i and a[i] == 'A':
                        first_A = q
                        break

                if a[5] == 'B':
                    valid = (x[5] == ((level <= qi[2]) or (first_A == 2)))
                elif a[5] == 'C':
                    valid = (x[5] == ((level <= qi[3]) or (first_A == 3)))
                elif a[5] == 'D':
                    valid = (x[5] == ((level <= qi[4]) or (first_A == 4)))
                elif a[5] == 'E':
                    valid = (x[5] == ((level <= qi[5]) or (first_A == 5)))



        # a[6] - Question 17
        i = qi[17]
        if valid and level > i:
            if DEBUG: print('q17')

            j = qi[10]

            if a[i] == 'E':
                valid = False
            elif level > j:
                if a[i] == 'A':
                    valid = (x[i] == (a[j] == 'C'))
                elif a[6] == 'B':
                    valid = (x[i] == (a[j] == 'D'))
                elif a[6] == 'C':
                    valid = (x[i] == (a[j] == 'B'))
                elif a[6] == 'D':
                    valid = (x[i] == (a[j] == 'A'))


        # a[7] - Question 10
        if valid and level > 7:
            if DEBUG: print('q10')

            if a[7] == 'E':
                valid = False
            elif level > qi[17]:
                if a[7] == 'A':
                    valid = (x[7] == (a[qi[17]] == 'D'))
                elif a[7] == 'B':
                    valid = (x[7] == (a[qi[17]] == 'B'))
                elif a[7] == 'C':
                    valid = (x[7] == (a[qi[17]] == 'A'))
                elif a[7] == 'D':
                    valid = (x[7] == (a[qi[17]] == 'E'))


        # a[8] - Question 5
        i, j = qi[5], qi[14]
        if valid and level > i and level > j:
            if DEBUG: print('q5')

            if a[i] == 'A':
                valid = (a[j] == 'B')
            if a[i] == 'B':
                valid = (a[j] == 'E')
            if a[i] == 'C':
                valid = (a[j] == 'C')
            if a[i] == 'D':
                valid = (a[j] == 'A')
            if a[i] == 'E':
                valid = (a[j] == 'D')

        # a[9] - Question 4
        i = qi[4]
        if valid and level > i:
            if DEBUG: print('q4')

            if a[i] == 'A' and level > qi[10] and level > qi[13]:
                valid = (x[i] == (a[qi[10]] == 'A' and a[qi[13]] == 'A'))

            elif a[i] == 'B' and level > qi[14] and level > qi[16]:
                valid = (x[i] == (a[qi[14]] == 'B' and a[qi[16]] == 'B'))

            elif a[i] == 'C' and level > qi[7] and level > qi[20]:
                valid = (x[i] == (a[qi[7]] == 'C' and a[qi[20]] == 'C'))

            elif a[i] == 'D' and level > qi[1] and level > qi[15]:
                valid = (x[i] == (a[qi[1]] == 'D' and a[qi[15]] == 'D'))

            elif a[i] == 'E' and level > qi[8] and level > qi[12]:
                valid = (x[i] == (a[qi[8]] == 'E' and a[qi[12]] == 'E'))

        # a[10] - Question 16
        i, j = qi[16], qi[8]
        if valid and level > i and level > j:
            if DEBUG: print('q16')

            if a[i] == 'A':
                k = qi[3]
            elif a[i] == 'B':
                k = qi[2]
            elif a[i] == 'C':
                k = qi[13]
            elif a[i] == 'D':
                k = qi[18]
            elif a[i] == 'E':
                k = qi[20]

            valid = (x[i] == ((k >= level) or (a[j] == a[k])))

        # a[11] - Question 11
        if valid and level > qi[11]:
            if DEBUG: print('q11')

            count_D = count['D']

            if a[qi[11]] == 'A':
                n = 2
            elif a[qi[11]] == 'B':
                n = 3
            elif a[qi[11]] == 'C':
                n = 4
            elif a[qi[11]] == 'D':
                n = 5
            elif a[qi[11]] == 'E':
                n = 6

            valid = (x[qi[11]] ==
                     ((count_D <= n) and ((level != 20) or (count_D == n))))


        # a[12] - Question 13
        i = qi[13]
        if valid and level > i:
            if DEBUG: print('q13')

            count_E = count['E']

            if a[i] == 'A':
                n = 5
            elif a[i] == 'B':
                n = 4
            elif a[i] == 'C':
                n = 3
            elif a[i] == 'D':
                n = 2
            elif a[i] == 'E':
                n = 1

            valid = (x[i] ==
                     ((count_E <= n) and ((level != 20) or (count_E == n))))

        # Question 14
        # TODO: is there a way to cutoff, before we reach l=20 ?
        i = qi[14]
        if valid and level == 20:
            if DEBUG: print('q14')

            if a[i] == 'A':
                valid = (x[i] == (2 not in values))
            elif a[i] == 'B':
                valid = (x[i] == (3 not in values))
            elif a[i] == 'C':
                valid = (x[i] == (4 not in values))
            elif a[i] == 'D':
                valid = (x[i] == (5 not in values))
            elif a[i] == 'E':
                valid = (x[i] == (values == [2, 3, 4, 5, 6]))


        # a[14] - Question 7
        i = qi[7]
        if valid and level == 20:
            if DEBUG: print('q7')

            most_often = [l for l in count.keys() if count[l] == values[-1]]

            if a[i] == 'A':
                valid = (x[i] == ('A' in most_often))
            elif a[i] == 'B':
                valid = (x[i] == ('B' in most_often))
            elif a[i] == 'C':
                valid = (x[i] == ('C' in most_often))
            elif a[i] == 'D':
                valid = (x[i] == ('D' in most_often))
            elif a[i] == 'E':
                valid = (x[i] == ('E' in most_often))

        # a[15] - Question 18
        i = qi[18]
        if valid and level > max_primes:
            if DEBUG: print('q18')

            vowels = 0
            for j in primes:
                if a[j] in ('A', 'E'):
                    vowels += 1

            if a[i] == 'A':
                valid = (x[i] == (vowels in (2, 3, 5, 7, 11, 13, 17, 19)))
            elif a[i] == 'B':
                valid = (x[i] == (vowels in (4, 9, 16)))
            elif a[i] == 'C':
                valid = (x[i] == (vowels % 2 == 1))
            elif a[i] == 'D':
                valid = (x[i] == (vowels % 2 == 0))
            elif a[i] == 'E':
                valid = (x[i] == (vowels == 0))

        # a[16] - Question 6
        i = qi[6]
        if valid and level > i:
            if DEBUG: print('q6')
            valid = (x[i] == True)

        # a[17] - Question 8
        i = qi[8]
        if valid and level == 20:
            # if DEBUG: print('q8')

            # if DEBUG: print(f'{values=}, {count=}')

            j = 0
            while j < len(values):
                min_count = values.count(values[j])
                if min_count == 1:
                    break
                j += min_count

            # if DEBUG: print(f'{min_count=}, {j=}')

            if j < len(values):
                # Get answer with this count
                for answer, c in count.items():
                    if c == values[j]:
                        break

                # if DEBUG: print(f'{x[i]=}, {a[i]=}, {answer=}')

                valid = (x[i] == (a[i] == answer))

        # a[18] - Question 12
        i = qi[12]
        if valid and level == 20:
            if DEBUG: print('q12')

            if DEBUG: print(f'{a[i]=}, {count=}')
            if a[i] == 'A':
                valid = (x[i] == ((count['A'] - 1) == count['B']))
            elif a[i] == 'B':
                valid = (x[i] == ((count['B'] - 1) == count['C']))
            elif a[i] == 'C':
                valid = (x[i] == ((count['C'] - 1) == count['D']))
            elif a[i] == 'D':
                valid = (x[i] == ((count['D'] - 1) == count['E']))
            elif a[i] == 'E':
                valid = (x[i] == (((count['A'] - 1) != count['B']) and
                                  ((count['B'] - 1) != count['C']) and
                                  ((count['C'] - 1) != count['D']) and
                                  ((count['D'] - 1) != count['E'])))


        # a[19] - Question 9
        i = qi[9]
        if valid and level == 20:
            if DEBUG: print('q9')

            # Sum of correct answers same as this one
            c = 0
            for j in range(20):
                if a[i] == a[j] and x[j]:
                    c += iq[j]

            if a[i] == 'A':
                valid = (x[i] == (59 <= c <= 62))
            elif a[i] == 'B':
                valid = (x[i] == (52 <= c <= 55))
            elif a[i] == 'C':
                valid = (x[i] == (44 <= c <= 49))
            elif a[i] == 'D':
                valid = (x[i] == (61 <= c <= 67))
            elif a[i] == 'E':
                valid = (x[i] == (44 <= c <= 53))


        if DEBUG:
            print(f'l={level}: {"Valid" if valid else "invalid"}: ' +
                        f'{to_str(level, a)}')

        return(valid)

    def to_str(level, a):
        ''' Convert answers to string representation. '''

        nonlocal x, iq

        s = []
        # for question in range(1, 21):
        #     i = q.index(question)
        #     if i < level:
        for i in range(level):
            s.append(str(iq[i]) + (a[i] if x[i] else a[i].lower()))

        return ' '.join(s)

    # Question map: index to question number
    iq = (3, 15, 20, 19, 2, 1, 17, 10, 5, 4, 16, 11, 13, 14, 7, 18, 6, 8, 12, 9)

    # Question map: question number to index
    qi = [0] * 21
    for i in range(20):
        qi[iq[i]] = i

    # Grading, in index order
    x = [None] * 20
    for q in range(1, 21):
        x[qi[q]] = grading[q-1]

    primes = [qi[2], qi[3], qi[5], qi[7], qi[11], qi[13], qi[17], qi[19]]
    max_primes = max(primes)

    for a in basic_backtrack(20, domain, property, stats=stats):
        # Convert index order to question order
        aq = []
        for q in range(1, 21):
            i = qi[q]
            aq.append(a[i] if x[i] else a[i].lower())

        yield tuple(aq)


class Test(unittest.TestCase):

    def test_exercise_71_warmup(self):
        result = list(exercise_71_warmup([False, False]))
        self.assertEqual(result, [('A', 'A'), ('B', 'A')])

        result = list(exercise_71_warmup([False, True]))
        self.assertEqual(result, [])

        result = list(exercise_71_warmup([True, False]))
        self.assertEqual(result, [('A', 'B')])

        result = list(exercise_71_warmup([True, True]))
        self.assertEqual(result, [('A', 'B')])

    def test_exercise_71_all_true(self):

        stats = {}
        grading = [True] * 20
        result = list(exercise_71(grading, stats))

        logger.info(f'all true: {sum(stats["level_count"])}, {stats}')

        self.assertEqual(result, [])

    def test_exercise_71_one_false(self):

        for q in (19, 20):
            stats = {}
            grading = [True] * 20
            grading[q-1] = False
            result = list(exercise_71(grading, stats))

            logger.info(f'{q} false: {sum(stats["level_count"])}, {stats}')

            if q == 19:
                self.assertEqual(result,
                                [('D', 'C', 'E', 'A', 'B', 'E', 'B', 'C', 'E',
                                  'A', 'B', 'E', 'A', 'E', 'D', 'B', 'D', 'A',
                                  'b', 'B')])

            elif q == 20:
                self.assertEqual(result,
                                 [('A', 'E', 'D', 'C', 'A', 'B', 'C', 'D', 'C',
                                   'A', 'C', 'E', 'D', 'B', 'C', 'A', 'D', 'A',
                                   'A', 'c'),
                                  ('D', 'C', 'E', 'A', 'B', 'A', 'D', 'C', 'D',
                                  'A', 'E', 'D', 'A', 'E', 'D', 'B', 'D', 'B',
                                  'E', 'e')])


if __name__ == '__main__':

    logger.addHandler(logging.StreamHandler())
    logger.setLevel(logging.INFO)

    unittest.main(exit=False)