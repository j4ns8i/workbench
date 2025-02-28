{{/*
Expand the name of the chart.
*/}}
{{- define "workbench.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "workbench.fullname" -}}
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

{{- define "workbench.apiName" -}}
{{ include "workbench.fullname" . }}-api
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "workbench.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "workbench.commonLabels" -}}
helm.sh/chart: {{ include "workbench.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Common selector labels
*/}}
{{- define "workbench.commonSelectorLabels" -}}
app.kubernetes.io/name: {{ include "workbench.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
API labels
*/}}
{{- define "workbench.apiLabels" -}}
{{ include "workbench.commonLabels" . }}
{{ include "workbench.commonSelectorLabels" . }}
{{ include "workbench.apiSelectorLabels" . }}
{{- end }}

{{/*
API Selector labels
*/}}
{{- define "workbench.apiSelectorLabels" -}}
workbench.j4ns8i.github.com/component: api
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "workbench.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "workbench.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Redis host
*/}}
{{- define "workbench.redisHost" -}}
{{- default "redis-master" .Values.redisHost }}
{{- end }}

{{/*
Redis host
*/}}
{{- define "workbench.redisPort" -}}
{{- default 6379 .Values.redisPort }}
{{- end }}

{{/*
Redis password Secret name
*/}}
{{- define "workbench.redisSecretName" -}}
{{- default "redis" .Values.redisPasswordName }}
{{- end }}

{{/*
Redis password Secret key
*/}}
{{- define "workbench.redisSecretPasswordKey" -}}
{{- default "redis-password" .Values.redisPasswordKey }}
{{- end }}
