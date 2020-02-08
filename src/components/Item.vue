<template>
  <div>
    <v-container class="mx-auto py-1 px-0">
      <v-card class="mx-auto" width="270" :color="item.color">
        <v-row justify="center" @click="openDialog">
          <!-- Item title -->
          <v-col cols="9">
            <v-card-title class="headline text-truncate py-1">{{ item.title }}</v-card-title>
            <v-card-subtitle class="py-2 text-truncate">{{ item.text }}</v-card-subtitle>
          </v-col>

          <!-- Menu button -->
          <v-col class="pt-0" cols="3">
            <div class="text-center">
              <v-menu offset-y>
                <template v-slot:activator="{ on }">
                  <v-btn fab outlined small text dark v-on="on">
                    <v-icon color="black">mdi-dots-horizontal</v-icon>
                  </v-btn>
                </template>
                <v-list>
                  <v-list-item
                    class="pr-1"
                    v-for="(item, index) in menuItems"
                    :key="index"
                    @click="item.action()"
                  >
                    <v-list-item-title>{{ item.title }}</v-list-item-title>
                    <v-list-item-avatar class="mx-0">
                      <v-icon small>{{ item.icon }}</v-icon>
                    </v-list-item-avatar>
                  </v-list-item>
                </v-list>
              </v-menu>
            </div>
          </v-col>
        </v-row>
      </v-card>
    </v-container>
    <Dialog
      v-model="editedItem"
      :open="editItemDialog"
      @close="editItemDialog = false"
      @save="saveEditedItem"
    ></Dialog>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import Dialog from '@/components/ItemModal.vue';

export default {
  name: 'item',
  components: {
    Dialog,
  },
  props: {
    id: String,
  },
  data() {
    return {
      editItemDialog: false,
      newItemTitle: '',
      editedItem: {},
      menuItems: [
        {
          title: 'edit',
          icon: 'mdi-pencil',
          action: this.edit,
        },
        {
          title: 'delete',
          icon: 'mdi-trash-can-outline',
          action: this.del,
        },
      ],
    };
  },
  computed: {
    item() {
      return this.getItemById()(this.id);
    },
  },
  methods: {
    ...mapGetters(['getItemById']),
    ...mapActions(['editItem', 'removeItem']),
    del() {
      this.removeItem(this.item);
    },
    edit() {
      this.newItemTitle = this.item.title;
      this.editItemDialog = true;
    },
    saveItemTitle() {
      const newTitle = this.newItemTitle.trim();
      if (newTitle) {
        const newItem = Object.assign({ ...this.item }, { title: newTitle });
        this.editItem(newItem);
        this.editItemDialog = false;
        this.newItemTitle = '';
      }
    },
    openDialog() {
      this.editedItem = { ...this.item };
      this.editItemDialog = true;
    },
    saveEditedItem() {
      this.editItem(this.editedItem);
      this.editItemDialog = false;
    },
  },
};
</script>
