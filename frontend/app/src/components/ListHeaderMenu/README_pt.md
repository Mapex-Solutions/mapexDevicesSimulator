# Componente ListHeaderMenu

> 🇺🇸 English version: [README.md](./README.md)

Botão de menu genérico para telas de lista que oferece contagem de itens, seleção de
itens por página e alternância de visibilidade de colunas.

## Funcionalidades

- **Contagem de itens**: mostra o total com rótulo no singular/plural
- **Itens por página**: dropdown para escolher o tamanho da paginação
- **Visibilidade de colunas**: checkboxes para alternar a exibição de colunas
- **Customizável**: configure ícone, rótulos e opções disponíveis
- **Two-way binding**: atualizações reativas ao componente pai

## Uso

```vue
<template>
  <ListHeaderMenu
    :items-count="assetsList.length"
    item-label="Asset"
    item-label-plural="Assets"
    icon="devices"
    :items-per-page="itemsPerPage"
    :columns="menuColumns"
    @update:items-per-page="itemsPerPage = $event"
    @update:columns="handleColumnsUpdate"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ListHeaderMenu } from '@components/headers';
import type { ListHeaderMenuColumn } from '@components/headers';

const itemsPerPage = ref(25);

const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'uuid', label: 'UUID', visible: true },
  { key: 'type', label: 'Type', visible: true },
  { key: 'protocol', label: 'Protocol', visible: true },
  { key: 'category', label: 'Category', visible: false },
]);

function handleColumnsUpdate(columns: ListHeaderMenuColumn[]) {
  menuColumns.value = columns;
  // Atualize aqui a sua lógica de colunas visíveis
}
</script>
```

## Props

| Prop | Tipo | Padrão | Descrição |
|------|------|--------|-----------|
| `itemsCount` | `number` | **obrigatório** | Número total de itens a exibir |
| `itemLabel` | `string` | **obrigatório** | Rótulo no singular (ex.: "Asset", "Rule", "User") |
| `itemLabelPlural` | `string` | `${itemLabel}s` | Rótulo no plural (ex.: "Assets", "Rules", "Users") |
| `icon` | `string` | `'list'` | Nome do ícone do botão |
| `itemsPerPage` | `number` | **obrigatório** | Valor atual de itens por página |
| `itemsPerPageOptions` | `number[]` | `[10, 25, 50, 100]` | Opções disponíveis de itens por página |
| `columns` | `ListHeaderMenuColumn[]` | `[]` | Configuração de visibilidade de colunas |
| `showItemsPerPage` | `boolean` | `true` | Mostra/oculta a seção de itens por página |
| `showColumnVisibility` | `boolean` | `true` | Mostra/oculta a seção de visibilidade de colunas |

## Configuração de colunas

### Interface ListHeaderMenuColumn

```typescript
interface ListHeaderMenuColumn {
  key: string;      // Identificador único da coluna
  label: string;    // Rótulo exibido no menu
  visible: boolean; // Estado atual de visibilidade
}
```

## Eventos

| Evento | Payload | Descrição |
|--------|---------|-----------|
| `update:itemsPerPage` | `value: number` | Emitido quando itens por página muda |
| `update:columns` | `columns: ListHeaderMenuColumn[]` | Emitido quando a visibilidade de colunas muda |

## Formato do rótulo do botão

O componente formata o rótulo do botão automaticamente:
- **Singular**: "1 ASSET"
- **Plural (padrão)**: "6 ASSETS"
- **Plural (customizado)**: "6 ITEMS" (se `itemLabelPlural="Items"`)

## Seções do menu

### Seção de itens por página

Aparece quando `showItemsPerPage !== false` e `itemsPerPageOptions.length > 0`.

- Mostra seleção em estilo de rádio
- Seleção atual marcada com check
- Fecha o menu na seleção
- Emite o evento `update:itemsPerPage`

### Seção de visibilidade de colunas

Aparece quando `showColumnVisibility !== false` e `columns.length > 0`.

- Mostra checkboxes para cada coluna
- Alterna colunas individualmente
- Atualizações em tempo real
- Emite o evento `update:columns`

## Exemplos

### Uso básico (só contagem de itens)

```vue
<ListHeaderMenu
  :items-count="items.length"
  item-label="Item"
  :items-per-page="25"
  :show-items-per-page="false"
  :show-column-visibility="false"
/>
```

### Só itens por página

```vue
<ListHeaderMenu
  :items-count="users.length"
  item-label="User"
  :items-per-page="itemsPerPage"
  :items-per-page-options="[5, 10, 20, 50]"
  :show-column-visibility="false"
  @update:items-per-page="itemsPerPage = $event"
/>
```

### Só visibilidade de colunas

```vue
<ListHeaderMenu
  :items-count="assets.length"
  item-label="Asset"
  :items-per-page="25"
  :columns="columns"
  :show-items-per-page="false"
  @update:columns="handleColumnsUpdate"
/>
```

### Configuração completa

```vue
<ListHeaderMenu
  :items-count="rules.length"
  item-label="Rule"
  item-label-plural="Rules"
  icon="rule"
  :items-per-page="itemsPerPage"
  :items-per-page-options="[10, 25, 50, 100]"
  :columns="menuColumns"
  @update:items-per-page="itemsPerPage = $event"
  @update:columns="handleColumnsUpdate"
/>
```

### Rótulos customizados

```vue
<ListHeaderMenu
  :items-count="12"
  item-label="Entry"
  item-label-plural="Entries"  <!-- "12 ENTRIES" em vez de "12 ENTRYS" -->
  icon="storage"
  :items-per-page="25"
/>
```

## Integração com DataRow

Use junto com o componente DataRow para uma lista completa:

```vue
<template>
  <!-- Cabeçalho com menu -->
  <div class="row items-center q-mb-md">
    <div class="col">
      <div class="text-subtitle1 text-weight-medium text-primary">
        Assets List
      </div>
    </div>
    <div class="col-auto">
      <ListHeaderMenu
        :items-count="filteredAssets.length"
        item-label="Asset"
        icon="devices"
        :items-per-page="itemsPerPage"
        :columns="menuColumns"
        @update:items-per-page="itemsPerPage = $event"
        @update:columns="handleColumnsUpdate"
      />
    </div>
  </div>

  <!-- Linhas de dados -->
  <div class="row">
    <div
      v-for="asset in paginatedAssets"
      :key="asset.id"
      class="col-12 q-mb-xs"
    >
      <DataRow
        :data="asset"
        :columns="visibleColumns"
        @edit="handleEdit"
        @view="handleView"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { ListHeaderMenu, DataRow } from '@components';
import type { ListHeaderMenuColumn, DataRowColumn } from '@components';

// Estado do menu
const itemsPerPage = ref(25);
const menuColumns = ref<ListHeaderMenuColumn[]>([
  { key: 'uuid', label: 'UUID', visible: true },
  { key: 'type', label: 'Type', visible: true },
  { key: 'protocol', label: 'Protocol', visible: true },
  { key: 'category', label: 'Category', visible: true },
]);

// Definições de colunas
const allColumns: DataRowColumn[] = [ /* ... */ ];

// Colunas visíveis computadas conforme a seleção do menu
const visibleColumns = computed(() => {
  return allColumns.filter(col => {
    const menuCol = menuColumns.value.find(mc => mc.key === col.key);
    return !menuCol || menuCol.visible;
  });
});

function handleColumnsUpdate(columns: ListHeaderMenuColumn[]) {
  menuColumns.value = columns;
}
</script>
```

## Estilização

O componente usa os estilos padrão de botão e menu do Quasar. Nenhum CSS extra é
necessário.

### Customização

Sobrescreva os estilos do botão se precisar:

```vue
<style scoped>
:deep(.q-btn) {
  font-weight: 600; /* Peso de fonte customizado */
}
</style>
```

## Acessibilidade

- Menu navegável por teclado
- Rótulos amigáveis a leitores de tela
- Atributos ARIA para checkboxes e itens de rádio
- Gerenciamento de foco

## Componentes relacionados

- **DataRow**: componente de linha para exibir itens de lista
- **PageHeader**: título da página e botões de ação
- **ListFilter**: componente de filtro para listas

## Casos de uso

Este componente é ideal para:
- Listas de assets
- Listas de regras
- Listas de usuários
- Listas de dispositivos
- Qualquer lista paginada com visibilidade de colunas opcional

## Suporte de navegador

Funciona em todos os navegadores modernos com suporte a:
- CSS Grid
- CSS Flexbox
- ES2022
