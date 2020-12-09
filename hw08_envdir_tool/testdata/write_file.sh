#!/usr/bin/env bash


TESTING_FILE=testing_written_file.txt

rm $TESTING_FILE 2> /dev/null
touch $TESTING_FILE
echo ${HELLO}
echo ${HELLO} >> "${TESTING_FILE}"
echo ${FOO}
echo ${FOO} >> "${TESTING_FILE}"
echo ${UNSET}
echo ${UNSET} >> "${TESTING_FILE}"
echo ${BAR}
echo ${BAR} >> "${TESTING_FILE}"