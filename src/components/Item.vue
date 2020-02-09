<template>
  <div>
    <v-card class="mx-auto my-1" width="270" @click.stop="edit">
      <v-container class="mx-auto py-1 px-0">
        <v-row v-if="item.tags.length" justify="start">
          <v-col cols="auto" class="py-0">
            <v-chip
              v-for="(tag, index) in tags"
              :key="index"
              :color="tag.color"
              class="ml-2"
              small
              label
              text-color="white"
            >
              <v-icon small left>mdi-label</v-icon>
              {{ tag.title }}
            </v-chip>
          </v-col>
        </v-row>
        <v-card-title class="headline text-truncate py-3">{{ item.title }}</v-card-title>
        <v-card-subtitle class="py-2 text-truncate">{{ item.text }}</v-card-subtitle>
      </v-container>
    </v-card>
    <ItemModal
      v-model="editedItem"
      :open="editItemMode"
      @close="saveEditedItem"
      @delete="del"
    ></ItemModal>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import ItemModal from '@/components/ItemModal.vue';

export default {
  name: 'item',
  components: {
    ItemModal,
  },
  props: {
    id: String,
  },
  data() {
    return {
      editItemMode: false,
      newItemTitle: '',
      editedItem: { title: '', text: '', tags: [] },
    };
  },
  computed: {
    ...mapGetters(['getItemById', 'getTagById']),
    item() {
      return this.getItemById(this.id);
    },
    tags() {
      return this.item.tags.map(tagId => this.getTagById(tagId));
    },
  },
  methods: {
    ...mapActions(['editItem', 'removeItem']),
    edit() {
      this.editedItem = Object.assign(
        { ...this.item },
        { tags: this.item.tags.slice() },
      );
      this.editItemMode = true;
    },
    saveItemTitle() {
      const newTitle = this.newItemTitle.trim();
      if (newTitle) {
        const newItem = Object.assign({ ...this.item }, { title: newTitle });
        this.editItem(newItem);
        this.editItemMode = false;
        this.newItemTitle = '';
      }
    },
    saveEditedItem() {
      this.editItem(this.editedItem);
      this.editItemMode = false;
    },
    del() {
      this.editItemMode = false;
      this.removeItem(this.item);
    },
  },
};
</script>
