# -*- coding: utf-8 -*- #
# Copyright 2020 Google LLC. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""The command group for gcloud dataflow sql."""

from __future__ import absolute_import
from __future__ import division
from __future__ import unicode_literals

from googlecloudsdk.calliope import base


@base.ReleaseTracks(base.ReleaseTrack.BETA, base.ReleaseTrack.GA)
@base.Deprecate(
    is_removed=False,
    warning=(
        'This command is deprecated and will be removed January 31, 2025. '
        'Please see [Beam YAML]'
        '(https://beam.apache.org/documentation/sdks/yaml/) '
        'and [Beam notebooks]'
        '(https://cloud.google.com/dataflow/docs/guides/notebook-advanced#beam-sql) '
        'for alternatives.'
    ),
    error=(
        'This command has been removed. '
        'Please see [Beam YAML]'
        '(https://beam.apache.org/documentation/sdks/yaml/) '
        'and [Beam notebooks]'
        '(https://cloud.google.com/dataflow/docs/guides/notebook-advanced#beam-sql) '
        'for alternatives.'
    ),
)
@base.DefaultUniverseOnly
class Sql(base.Group):
  """A group of subcommands for working with Dataflow SQL."""

  pass