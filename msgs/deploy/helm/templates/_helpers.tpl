{{/*
Expand the name of the chart.
*/}}
{{- define "msgs.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "msgs.fullname" -}}
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
{{- define "msgs.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "msgs.commonLabels" -}}
helm.sh/chart: {{ include "msgs.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Labels
*/}}
{{- define "msgs.labels" -}}
{{ include "msgs.commonLabels" . }}
{{ include "msgs.selectorLabels" . }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "msgs.selectorLabels" -}}
app.kubernetes.io/name: {{ include "msgs.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
workbench.j4ns8i.github.com/component: msgs
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "msgs.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "msgs.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Redis host
*/}}
{{- define "msgs.redisHost" -}}
{{- default "redis-master" .Values.redisHost }}
{{- end }}

{{/*
Redis host
*/}}
{{- define "msgs.redisPort" -}}
{{- default 6379 .Values.redisPort }}
{{- end }}

{{/*
Redis password Secret name
*/}}
{{- define "msgs.redisSecretName" -}}
{{- default "redis" .Values.redisPasswordName }}
{{- end }}

{{/*
Redis password Secret key
*/}}
{{- define "msgs.redisSecretPasswordKey" -}}
{{- default "redis-password" .Values.redisPasswordKey }}
{{- end }}
