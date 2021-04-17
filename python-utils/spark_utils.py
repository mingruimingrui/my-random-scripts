from __future__ import absolute_import, unicode_literals, print_function

import os
import subprocess

from pyspark import SparkConf, SparkContext
from pyspark.sql import SparkSession, DataFrame


def get_spark() -> SparkSession:
    SparkSession.builder.enableHiveSupport().getOrCreate()


def create_spark(
    app_name: str = 'mingrui_playground',
    **extra_kwargs
) -> SparkSession:
    """Create spark context and return session."""
    kwargs = {
        'spark.driver.allowMultipleContexts': True,
        'spark.locality.wait': '500ms',
        'spark.executor.memory': '1g',
        'spark.driver.maxResultSize': '32g',
        'spark.executor.instances': '16',
        'spark.executor.cores': '3',
        'spark.memory.fraction': '0.2',
        'spark.sql.crossJoin.enabled': True,
    }  # Default kwargs
    for k, v in extra_kwargs.items():
        k = k.replace('_', '.').replace('-', '.')
        kwargs[k] = v

    conf = SparkConf().setAppName(app_name)
    for k, v in kwargs.items():
        conf = conf.set(k, v)
    SparkContext(conf=conf)

    return SparkSession.builder.enableHiveSupport().getOrCreate()


def df_to_csv(
    df: DataFrame,
    local_filepath: str,
    sep: str = ',',
    header: bool = True
):
    """Write a spark dataframe to a local CSV file."""
    local_filepath = os.path.abspath(local_filepath)
    local_dirname = os.path.dirname(local_filepath)
    local_basename = os.path.basename(local_filepath)

    hdfs_filepath = 'tmp/{0}'.format(local_basename)
    crc_filepath = '{}/.{}.crc'.format(local_dirname, local_basename)

    # Create context to handle cleanup
    class DFToCSVContext():

        def __enter__(self):
            return self

        def __exit__(self, *exc_args):
            if os.path.isfile(crc_filepath):
                os.remove(crc_filepath)
            subprocess.call(
                'hdfs dfs -rm -r {}'.format(hdfs_filepath),
                shell=True
            )
            return True

    with DFToCSVContext():
        # Write to hdfs
        print('Writing to hdfs')
        df.write.csv(
            path=hdfs_filepath,
            mode='overwrite',
            sep=sep,
            header=header
        )

        # Get merge and remove hdfs file
        print('Copying to local')
        subprocess.call('hdfs dfs -getmerge {} {}'.format(
            hdfs_filepath, local_filepath
        ), shell=True)
