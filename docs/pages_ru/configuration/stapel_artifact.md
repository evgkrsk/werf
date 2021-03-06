---
title: Stapel-артефакт
sidebar: documentation
permalink: documentation/configuration/stapel_artifact.html
author: Alexey Igrychev <alexey.igrychev@flant.com>
summary: |
  <a class="google-drawings" href="../../images/configuration/stapel_artifact1.png" data-featherlight="image">
      <img src="../../images/configuration/stapel_artifact1_preview.png">
  </a>
---

## Что такое артефакты?

***Артефакт*** — это специальный образ, используемый в других артефактах или отдельных образах, описанных в конфигурации. Артефакт предназначен преимущественно для отделения ресурсов инструментов сборки от процесса сборки образа приложения. Примерами таких ресурсов могут быть — программное обеспечение или данные, которые необходимы для сборки, но не нужны для запуска приложения, и т.п.

Образ _артефакта_ нельзя [протегировать]({{ site.baseurl }}/documentation/reference/publish_process.html) как обычный образ, и использовать как отдельное приложение.

Используя артефакты, вы можете собирать неограниченное количество компонентов, что позволяет решать, например, следующие задачи:
- Если приложение состоит из набора компонент, каждый со своими зависимостями, то обычно вам приходится пересобирать все компоненты каждый раз. Вам бы хотелось пересобирать только те компоненты, которым это действительно нужно.
- Компоненты должны быть собраны в разных окружениях.

Импортирование _ресурсов_ из _артефактов_ указывается с помощью [директивы import]({{ site.baseurl }}/documentation/configuration/stapel_image/import_directive.html) в конфигурации в [_секции образа_]({{ site.baseurl }}/documentation/configuration/introduction.html#секция-образа) или [_секции артефакта_]({{ site.baseurl }}/documentation/configuration/introduction.html#секция-артефакта)).

## Конфигурация

Конфигурация _артефакта_ похожа на конфигурацию обычного _образа_. Каждый _артефакт_ должен быть описан в своей [секции]({{ site.baseurl }}/documentation/configuration/introduction.html#секция-артефакта) конфигурации.

Инструкции, связанные со стадией _from_ (инструкции указания [базового образа]({{ site.baseurl }}/documentation/configuration/stapel_image/base_image.html) и [монтирования]({{ site.baseurl }}/documentation/configuration/stapel_image/mount_directive.html)), а также инструкции [импорта]({{ site.baseurl }}/documentation/configuration/stapel_image/import_directive.html) точно такие же как и при описании _образа_.

Стадия добавления инструкций Docker (`docker_instructions`) и [соответствующие директивы]({{ site.baseurl }}/documentation/configuration/stapel_image/docker_directive.html) не доступны при описании _артефактов_. _Артефакт_ — это инструмент сборки, и все что от него требуется, это — только данные.

Остальные _стадии_ и инструкции описания артефактов рассматриваются далее подробно.

### Именование

<div class="summary" markdown="1">
```yaml
artifact: <artifact name>
```
</div>

_Образ артефакта_ объявляется с помощью директивы `artifact`. Синтаксис: `artifact: <artifact name>`. Так как артефакты используются только самим werf, отсутствуют какие-либо ограничения на именование артефактов, в отличие от ограничений на [именование обычных _образов_]({{ site.baseurl }}/documentation/configuration/stapel_image/naming.html).

Пример:
```yaml
artifact: "application assets"
```

### Добавление исходного кода из git-репозиториев

<div class="summary">

<a class="google-drawings" href="../../images/configuration/stapel_artifact2.png" data-featherlight="image">
  <img src="../../images/configuration/stapel_artifact2_preview.png">
</a>

</div>

В отличие от обычных _образов_, у _конвейера стадий артефактов_ нет стадий _gitCache_ и _gitLatestPatch_.

> В werf для _артефактов_ реализована необязательная зависимость от изменений в git-репозиториях. Таким образом, по умолчанию werf игнорирует какие-либо изменения в git-репозитории, кэшируя образ после первой сборки. Но вы можете определить зависимости от файлов и папок, при изменении в которых образ артефакта будет пересобираться

Читайте подробнее про работу с _git-репозиториями_ в соответствующей [статье]({{ site.baseurl }}/documentation/configuration/stapel_image/git_directive.html).

### Запуск инструкций сборки

<div class="summary">

<a class="google-drawings" href="../../images/configuration/stapel_artifact3.png" data-featherlight="image">
  <img src="../../images/configuration/stapel_artifact3_preview.png">
</a>

</div>

У артефактов точно такое же как и у обычных образов использование директив и пользовательских стадий — _beforeInstall_, _install_, _beforeSetup_ и _setup_.

Если в директиве `stageDependencies` в блоке git для _пользовательской стадии_ не указана зависимость от каких-либо файлов, то образ кэшируется после первой сборки, и не будет повторно собираться пока соответствующая _стадия_ находится в _stages storage_.

> Если необходимо повторно собирать артефакт при любых изменениях в git, нужно указать _stageDependency_ `**/*` для соответствующей _пользовательской_ стадии. Пример для стадии _install_:
```yaml
git:
- to: /
  stageDependencies:
    install: "**/*"
```

Читайте подробнее про работу с _инструкциями сборки_ в соответствующей [статье]({{ site.baseurl }}/documentation/configuration/stapel_image/assembly_instructions.html).

## Все директивы

{% include /configuration/stapel_image/artifact_directives.md %}

## Использование артефактов

В отличие от [*обычного образа*]({{ site.baseurl }}/documentation/configuration/stapel_image/assembly_instructions.html), у *образа артефакта* нет стадии _git latest patch_. Это сделано намеренно, т.к. стадия _git latest patch_ выполняется обычно при каждом коммите, применяя появившиеся изменения к файлам. Однако *артефакт* рекомендуется использовать как образ с высокой вероятностью кэширования, который обновляется редко или нечасто (например, при изменении специальных файлов).

Пример: нужно импортировать в артефакт данные из git, и пересобирать ассеты только тогда, когда влияющие на сборку ассетов файлы изменяются. Т.е. в случае, изменения каких-либо других файлов в git, ассеты пересобираться не должны.

Конечно, существуют случаи когда необходимо включать изменения любых файлов git-репозитория в _образ артефакта_ (например, если в артефакте происходит сборка приложения на Go). В этом случае необходимо указать зависимость относительно стадии (сборку которой необходимо выполнять при изменениях в git) с помощью `git.stageDependencies` и `*` в качестве шаблона. Пример:

```yaml
git:
- add: /
  to: /app
  stageDependencies:
    setup:
    - "*"
```

В этом случае любые изменения файлов в git-репозитории будут приводить к пересборке _образа артефакта_, и всех _образов_, в которых определен импорт этого артефакта.

**Замечание:** Если вы используете какие-либо файлы и при сборке _артефакта_ и при сборке [*обычного образа*]({{ site.baseurl }}/documentation/configuration/stapel_image/assembly_instructions.html), правильный путь — использовать директиву `git.add` при описании каждого образа, где это необходимо, т.е. несколько раз. **Не рекомендуемый** вариант — добавить файлы при сборке артефакта, а потом импортировать их используя директиву `import` в другой образ.
