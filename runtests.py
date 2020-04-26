import time
from unittest import TestCase, TestLoader, TextTestRunner, TestSuite
from unittest.runner import TextTestResult
import argparse

# Enhanced timing information
# https://hackernoon.com/timing-tests-in-python-for-fun-and-profit-1663144571


class TimeLoggingTestResult(TextTestResult):

    def __init__(self, stream, descriptions, verbosity):
        super().__init__(stream, descriptions, verbosity)
        self.test_timings = []
        self.stream = stream
        self.showAll = verbosity > 1

    def startTest(self, test):
        self._test_started_at = time.time()
        super().startTest(test)

    def addSuccess(self, test):
        elapsed = time.time() - self._test_started_at
        name = self.getDescription(test)
        self.test_timings.append((name, elapsed))

        if self.showAll:
            self.stream.write("({:.03}s)  ".format(round(elapsed, 3)))

        super().addSuccess(test)

    def getTestTimings(self):
        return self.test_timings


class TimeLoggingTestRunner(TextTestRunner):

    def __init__(self, slow_test_threshold=None, *args, **kwargs):
        self.slow_test_threshold = slow_test_threshold
        return super().__init__(
            resultclass=TimeLoggingTestResult,
            *args,
            **kwargs,
        )

    def run(self, test):
        result = super().run(test)

        results = []
        for name, elapsed in result.getTestTimings():
            if self.slow_test_threshold and \
               elapsed > self.slow_test_threshold:
                results.append((name, elapsed))

        if len(results) > 0 and self.slow_test_threshold is not None:
            s = "\nSlow Tests (>{:.03}s):".format(self.slow_test_threshold)
            self.stream.writeln(s)

            results.sort(key=lambda x: x[1])

            for name, elapsed in results:
                self.stream.writeln(
                    "({:.03}s) {}".format(
                        elapsed, name))

        return result


if __name__ == '__main__':
    # Get command line arguments
    parser = argparse.ArgumentParser()

    parser.add_argument("-v", "--verbosity",
                        type=int, default=1, choices=[0, 1, 2],
                        help="verbosity of test output; default is 1")

    parser.add_argument("-p", "--pattern",
                        default="*.py",
                        help="pattern of modules to test; default is *.py",)

    parser.add_argument('-d', '--directory',
                        help='root directory of the tests')

    parser.add_argument('-l', '--longtests', action='store_true',
                        help='run long running tests; default is false')

    parser.add_argument('-t', '--threshold',
                        type=float, default=5.0,
                        help=('threshold in seconds for considering a test '
                              'long running'))

    args = parser.parse_args()

    # Get list of test cases to run
    test_loader = TestLoader()
    test_suite_discovered = test_loader.discover(args.directory,
                                                 pattern=args.pattern)

    # Strip out long running tests, if necessary
    if args.longtests:
        test_suite = test_suite_discovered

    else:
        def iterate_tests(test_suite_or_case):
            """Iterate through all of the test cases in 'test_suite_or_case'.
            https://stackoverflow.com/questions/15487587/
            python-unittest-get-testcase-ids-from-nested-testsuite
            """
            try:
                suite = iter(test_suite_or_case)
            except TypeError:
                yield test_suite_or_case
            else:
                for test in suite:
                    for subtest in iterate_tests(test):
                        yield subtest

        test_suite = TestSuite()
        for test in iterate_tests(test_suite_discovered):
            if not test._testMethodName.startswith('test_long_'):
                test_suite.addTest(test)

    sth = None if args.longtests else args.threshold
    test_runner = TimeLoggingTestRunner(slow_test_threshold=sth,
                                        verbosity=args.verbosity)

    test_runner.run(test_suite)
