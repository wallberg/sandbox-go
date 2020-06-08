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


def exercise_71_warmup(n_correct_answers, stats=None):
    '''
    Generate solutions to Exercise 71 warmup (Two Questions) with n correct
    answers.
    '''

    assert 0 <= n_correct_answers <= 2

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

        nonlocal x, n_incorrect_answers

        # a[0] - Question 1
        if level > 0:
            if a[0] == 'B':
                x[0] = False
            elif a[0] == 'A' and level > 1:
                x[0] = (a[1] == 'B')
            else:
                x[0] = True

        # a[1] - Question 2
        if level > 1:
            if a[1] == 'A':
                x[1] = x[0]
            elif a[0] == 'B' and a[1] == 'B':
                # No valid grading
                return False
            else:
                x[1] = True

        if x[:level].count(False) > n_incorrect_answers:
            return False

        return True

    n = 2
    n_incorrect_answers = n - n_correct_answers

    # correctness of each answer
    x = [None] * n

    for a in basic_backtrack(n, domain, property, stats=stats):
        if x.count(True) == n_correct_answers:
            yield (a, tuple(x))


def exercise_71(n_correct_answers, stats=None):
    '''
    Generate solutions to Exercise 71 (Twenty Questions) with n correct
    answers.
    '''

    assert 0 <= n_correct_answers <= 20

    def domain(n, k):
        return ('A', 'B', 'C', 'D', 'E')

    def property(n, level, a):
        '''
        Check for count of correct answers, in this order of questions:
        3, 15, 20, 19, 2, 1, 17, 10, 5, 4, 16, 11, 13, 14, 7, 18, 6, 8, 12, 9
        '''

        nonlocal x

        # Question 3 - a[0]
        if level >= 1:
            x[0] = True
            for i in range(level-1):
                if a[i] == a[i+1]:
                    if a[0] == 'A' and i != 15:
                        x[0] = False
                        break
                    elif a[0] == 'B' and i != 16:
                        x[0] = False
                        break
                    elif a[0] == 'C' and i != 17:
                        x[0] = False
                        break
                    elif a[0] == 'D' and i != 18:
                        x[0] = False
                        break
                    elif a[0] == 'E' and i != 19:
                        x[0] = False
                        break

        # Question 15 - a[1]
        if level >= 2:
            x[1] = True

        return True

    n = 2
    n_correct_answers = 2
    n_incorrect_answers = 2 - n_correct_answers

    # correctness of each answer
    x = [None] * n

    for answers in basic_backtrack(n, domain, property, stats=stats):
        print(answers, x)
        yield answers


class Test(unittest.TestCase):

    def test_exercise_71_warmup(self):
        result = list(exercise_71_warmup(0))
        self.assertEqual(result, [(('A', 'A'), (False, False)),
                                  (('B', 'A'), (False, False))])

        result = list(exercise_71_warmup(1))
        self.assertEqual(len(result), 0)

        result = list(exercise_71_warmup(2))
        self.assertEqual(result, [(('A', 'B'), (True, True))])


if __name__ == '__main__':
    unittest.main(exit=False)

    logger.addHandler(logging.StreamHandler())
    logger.setLevel(logging.INFO)

    for n_correct_answers in range(3):
        stats = {}
        for a, x in exercise_71_warmup(n_correct_answers, stats): 
            print(f'correct={n_correct_answers}, {a=}, {x=}')
        print(f'correct={n_correct_answers}, {sum(stats["level_count"])}, {stats}')