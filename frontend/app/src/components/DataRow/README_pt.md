# Componente DataRow

> 🇺🇸 English version: [README.md](./README.md)

Componente de linha genérico para exibir dados em telas de lista com layouts responsivos.

## Funcionalidades

- **Layout Desktop/Laptop (>= 600px)**: cards em estilo de tabela com todas as colunas visíveis
- **Layout Mobile (< 600px)**: cards compactos com detalhes expansíveis
- **Breakpoints responsivos**: ocultação automática de colunas conforme o tamanho da tela
- **Menu de ações**: ações nativas de Editar, Ver e Excluir
- **Tipos de coluna**: Avatar, Text, Code, Chip, Badge
- **Texto secundário**: suporte a exibir texto adicional abaixo do conteúdo principal
- **Handlers de clique**: eventos de clique simples, duplo clique e ações

## Uso

```vue
<template>
  <DataRow
    :data="asset"
    :columns="assetColumns"
    primary-key="id"
    @click="handleClick"
    @dblclick="handleEdit"
    @edit="handleEdit"
    @view="handleView"
    @delete="handleDelete"
  />
</template>

<script setup lang="ts">
import { DataRow } from '@components/cards';
import type { DataRowColumn } from '@components/cards';

const assetColumns: DataRowColumn[] = [
  {
    key: 'icon',
    label: '',
    type: 'avatar',
    visible: 'always',
    width: 56,
    icon: (value, row) => row.icon || 'sensors',
    color: (value, row) => row.status ? 'primary' : 'grey-5',
  },
  {
    key: 'name',
    label: 'Name',
    type: 'text',
    visible: 'always',
    width: 250,
    ellipsis: true,
    secondaryKey: 'description', // Mostra a descrição abaixo do nome
  },
  {
    key: 'uuid',
    label: 'UUID',
    type: 'code',
    visible: 'always',
    width: 180,
    ellipsis: true,
  },
  {
    key: 'type',
    label: 'Type',
    type: 'chip',
    visible: 'always',
    width: 150,
    color: 'blue-6',
  },
  {
    key: 'protocol.type',
    label: 'Protocol',
    type: 'chip',
    visible: 'desktop', // Oculto em telas <= 1024px
    width: 120,
    format: (value) => value?.toUpperCase() || 'N/A',
  },
  {
    key: 'status',
    label: 'Status',
    type: 'badge',
    visible: 'always',
    width: 100,
    format: (value) => value ? 'ACTIVE' : 'INACTIVE',
    color: (value) => value ? 'green-6' : 'red-6',
  },
];
</script>
```

## Props

| Prop | Tipo | Padrão | Descrição |
|------|------|--------|-----------|
| `data` | `any` | **obrigatório** | O objeto de dados a exibir |
| `columns` | `DataRowColumn[]` | **obrigatório** | Array de configuração das colunas |
| `primaryKey` | `string` | `'id'` | Chave única da linha |
| `showActions` | `boolean` | `true` | Mostra/oculta o menu de ações |
| `expandOnClick` | `boolean` | `true` | Habilita a expansão do card mobile ao clicar |

## Configuração de colunas

### Interface DataRowColumn

```typescript
interface DataRowColumn {
  key: string;                    // Caminho da propriedade (suporta aninhado: 'protocol.type')
  label: string;                  // Rótulo do cabeçalho da coluna
  type: DataRowColumnType;        // Tipo: 'avatar' | 'text' | 'code' | 'chip' | 'badge'
  visible: DataRowColumnVisibility; // Visibilidade: 'always' | 'desktop' | 'laptop' | 'expandable'
  width?: number;                 // Largura da coluna em pixels
  ellipsis?: boolean;             // Habilita truncamento de texto com tooltip
  secondaryKey?: string;          // Chave do texto secundário (abaixo do texto principal)

  // Formatação e estilo
  format?: (value: any, row: any) => string;           // Formatador de valor
  color?: string | ((value: any, row: any) => string); // Cor (valor ou função)
  icon?: string | ((value: any, row: any) => string);  // Ícone (nome ou função, para avatar)
}
```

### Tipos de coluna

- **`avatar`**: exibe um ícone com cor de fundo
- **`text`**: texto simples com texto secundário opcional abaixo
- **`code`**: texto monoespaçado para UUIDs, IDs etc.
- **`chip`**: chip/tag colorido
- **`badge`**: badge colorido para indicadores de status

### Opções de visibilidade

- **`always`**: sempre visível em qualquer tamanho de tela
- **`desktop`**: oculto em telas <= 1024px
- **`laptop`**: oculto em telas <= 1366px
- **`expandable`**: só aparece na área expandida do mobile

## Eventos

| Evento | Payload | Descrição |
|--------|---------|-----------|
| `click` | `data: any` | Emitido ao clicar na linha |
| `dblclick` | `data: any` | Emitido no duplo clique da linha |
| `edit` | `data: any` | Emitido ao clicar na ação Editar |
| `view` | `data: any` | Emitido ao clicar na ação Ver |
| `delete` | `data: any` | Emitido ao clicar na ação Excluir |
| `expand` | `data: any, expanded: boolean` | Emitido ao expandir/recolher o card mobile |

## Comportamento responsivo

### Desktop/Laptop (>= 600px)
- Exibe como linha horizontal com todas as colunas visíveis
- Efeito de hover com sombra
- Menu de ações alinhado à borda direita
- Tooltips no texto truncado

### Mobile (< 600px)
- Card compacto mostrando: Avatar, Name+Description, Status, Ações
- Clique em qualquer lugar do card para expandir/recolher
- A área expandida mostra as colunas restantes em grid de 2 colunas
- Transição de slide suave

## Exemplos

### Uso básico

```vue
<DataRow
  :data="{ id: 1, name: 'Temperature Sensor', status: true }"
  :columns="columns"
/>
```

### Com todos os eventos

```vue
<DataRow
  :data="item"
  :columns="columns"
  @click="console.log('Clicked:', $event)"
  @dblclick="openEditor($event)"
  @edit="openEditor($event)"
  @view="openDetails($event)"
  @delete="confirmDelete($event)"
/>
```

### Acesso a propriedade aninhada

```vue
const columns = [
  {
    key: 'protocol.type',        // Acessa propriedade aninhada
    label: 'Protocol',
    type: 'chip',
    format: (value) => value?.toUpperCase(),
  },
];
```

### Estilo dinâmico

```vue
const columns = [
  {
    key: 'status',
    label: 'Status',
    type: 'badge',
    format: (value) => value ? 'ACTIVE' : 'INACTIVE',
    color: (value) => value ? 'green-6' : 'red-6', // Cor dinâmica
  },
];
```

## Estilização

O componente usa SCSS com escopo e classes utilitárias do Quasar. Todos os estilos são
aplicados automaticamente — nenhum CSS extra é necessário.

### Customização

Você pode sobrescrever estilos usando deep selectors:

```vue
<style scoped>
:deep(.data-row-card) {
  border-radius: 8px; /* Raio de borda customizado */
}

:deep(.data-row-card:hover) {
  background-color: #f0f0f0; /* Cor de hover customizada */
}
</style>
```

## Componentes relacionados

- **ListHeaderMenu**: menu de itens por página e visibilidade de colunas
- **Column Components**: AvatarColumn, TextColumn, CodeColumn, ChipColumn, BadgeColumn

## Suporte de navegador

Funciona em todos os navegadores modernos com suporte a:
- CSS Grid
- CSS Flexbox
- ES2022
- Media Queries
