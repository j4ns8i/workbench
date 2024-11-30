{{/*
Expand the name of the chart.
*/}}
{{- define "tilt.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tilt.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "tilt.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tilt.commonLabels" -}}
helm.sh/chart: {{ include "tilt.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Common selector labels
*/}}
{{- define "tilt.commonSelectorLabels" -}}
app.kubernetes.io/name: {{ include "tilt.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
API labels
*/}}
{{- define "tilt.apiLabels" -}}
{{ include "tilt.commonLabels" . }}
{{ include "tilt.commonSelectorLabels" . }}
{{ include "tilt.apiSelectorLabels" . }}
{{- end }}

{{/*
API Selector labels
*/}}
{{- define "tilt.apiSelectorLabels" -}}
workbench.j4ns8i.github.com/component: api
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tilt.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tilt.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
