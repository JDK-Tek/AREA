<template>
    <v-menu>
        <template v-slot:activator="{ props }">
            <v-btn
            color="primary"
            v-bind="props"
            >
            Catégorie
            </v-btn>
        </template>
        <v-list>
            <v-list-item
            v-for="(item, index) in items"
            :key="index"
            :value="index"
            >
            <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script>
    export default {
        data() {
            return {
                items: [
                    { title: 'Web' },
                    { title: 'Database' },
                    { title: 'Backend' },
                    { title: 'Mobile' },
                ]
            }
        }
    }
</script>