<template>
  <div>
    <v-card class="mx-auto" width="300" color="grey lighten-4">
      <!-- PIN: List -->
      <v-container>
        <v-row>
          <!-- PIN: List title -->
          <v-col cols="9">
            <v-card-title v-if="editListMode" color="grey lighten-4">
              <v-text-field
                class="headline ml-0 mt-0"
                v-model="newListTitle"
                dense
                full-width
                label="New List"
                single-line
                autofocus
                :rules="[!!newListTitle || 'Required']"
                @keypress.enter="saveListTitle()"
                @blur="saveListTitle()"
              ></v-text-field>
            </v-card-title>
            <v-card-title v-else class="headline">{{ list.title }}</v-card-title>
          </v-col>

          <!-- PIN: Menu button -->
          <v-col cols="3">
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
      </v-container>

      <v-divider class="mx-3"></v-divider>

      <!-- PIN: Items -->
      <v-list dense color="grey lighten-4">
        <draggable v-model="items" group="items">
          <v-list-item v-for="(item, index) in items" :key="index">
            <CardItem :id="item.id" />
          </v-list-item>
          <v-list-item v-if="addItemMode" class="my-2">
            <v-card class="mx-auto" width="270">
              <v-text-field
                class="headline mx-4"
                v-model="newItemTitle"
                dense
                full-width
                label="New Item"
                single-line
                autofocus
                :rules="[!!newItemTitle || 'Required']"
                @keypress.enter="addNewItem"
                @blur="addNewItem"
              ></v-text-field>
            </v-card>
          </v-list-item>
        </draggable>
      </v-list>

      <!-- PIN: Button for creating new item -->
      <v-card-actions>
        <v-btn
          class="mx-auto"
          width="250"
          depressed
          color="grey lighten-2"
          @click="addItemMode = true"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-card-actions>
    </v-card>

    <!-- PIN: Confirmation modal for deleting board -->
    <ModalConfirm
      :title="`Delete '${list.title}' ?`"
      text="Can not restore list and items in list."
      :open="deleteListMode"
      @cancel="deleteListMode = false"
      @confirm="del"
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import draggable from 'vuedraggable';
import ModalConfirm from '@/components/ModalConfirm.vue';
import CardItem from '@/components/CardItem.vue';

export default {
  name: 'CardList',
  components: {
    ModalConfirm,
    CardItem,
    draggable,
  },
  props: {
    id: String,
  },
  computed: {
    ...mapGetters(['getListById', 'getItemsByListId', 'getItemById']),
    list() {
      return this.getListById(this.id);
    },
    items: {
      get() {
        return this.getItemsByListId(this.id);
      },
      set(newValue) {
        const toList = Object.assign(
          { ...this.list },
          { items: newValue.map(item => item.id) },
        );
        if (this.items.length === toList.items.length) {
          // It Update this list if it moves item in this list.
          this.editList(toList);
        } else if (this.items.length < toList.items.length) {
          // This process runs if it moves item to this list.
          const newItemId = toList.items.find(newId => !this.list.items.includes(newId));
          const item = this.getItemById(newItemId);
          this.moveItemAcrossLists({ item, toList });
        }
      },
    },
  },
  data() {
    return {
      editListMode: false,
      deleteListMode: false,
      addItemMode: false,
      newListTitle: '',
      newItemTitle: '',
      menuItems: [
        {
          title: 'edit',
          icon: 'mdi-pencil',
          action: () => {
            this.newListTitle = this.list.title;
            this.editListMode = true;
          },
        },
        {
          title: 'delete',
          icon: 'mdi-trash-can-outline',
          action: () => {
            this.deleteListMode = true;
          },
        },
      ],
    };
  },
  methods: {
    ...mapActions(['editList', 'removeList', 'addItem', 'editItem', 'moveItemAcrossLists']),
    addNewItem() {
      const newTitle = this.newItemTitle.trim();
      if (newTitle) {
        this.addItem({
          listId: this.id,
          title: this.newItemTitle.trim(),
        });
      }
      this.newItemTitle = '';
      this.addItemMode = false;
    },
    del() {
      this.deleteListMode = false;
      this.removeList(this.list);
    },
    saveListTitle() {
      const newTitle = this.newListTitle.trim();
      if (newTitle) {
        const newList = Object.assign({ ...this.list }, { title: newTitle });
        this.editList(newList);
        this.editListMode = false;
        this.newListTitle = '';
      }
    },
  },
};
</script>
