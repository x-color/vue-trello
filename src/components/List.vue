<template>
  <div>
    <v-card class="mx-auto" width="300">
      <v-row>
        <!-- List title -->
        <v-col cols="9">
          <v-card-title
            v-if="!editListMode"
            class="headline text-truncate"
          >{{ list.title }}</v-card-title>
          <v-text-field
            class="headline ml-4 mt-3"
            v-model="newListTitle"
            v-if="editListMode"
            dense
            full-width
            label="New List"
            single-line
            autofocus
            :rules="[!!newListTitle || 'Required']"
            @keypress.enter="saveListTitle()"
            @blur="saveListTitle()"
          ></v-text-field>
        </v-col>

        <!-- Menu button -->
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

      <v-divider class="mx-3"></v-divider>

      <v-list>
        <!-- <v-list-item>
          <Item />
        </v-list-item>-->
      </v-list>
      <v-card-actions>
        <v-btn class="mx-auto" width="250" depressed @click="modal = true">
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-card-actions>
    </v-card>

    <ConfirmModal
      :title="`Delete '${list.title}' ?`"
      text="Can not restore list and items in list."
      :open="deleteListMode"
      @cancel="deleteListMode = false"
      @confirm="del()"
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import ConfirmModal from '@/components/ConfirmModal.vue';

export default {
  name: 'list',
  components: {
    ConfirmModal,
  },
  props: {
    id: String,
  },
  data() {
    return {
      modal: false,
      editListMode: false,
      deleteListMode: false,
      newListTitle: false,
      menuItems: [
        {
          title: 'edit',
          icon: 'mdi-pencil',
          action: this.edit,
        },
        {
          title: 'delete',
          icon: 'mdi-trash-can-outline',
          action: () => { this.deleteListMode = true; },
        },
      ],
    };
  },
  computed: {
    list() {
      return this.getListById()(this.id);
    },
  },
  methods: {
    ...mapGetters(['getListById']),
    ...mapActions(['editList', 'removeList']),
    closeDialog() {
      this.modal = false;
    },
    saveDialog() {
      this.modal = false;
    },
    del() {
      this.removeList(this.list);
      this.deleteListMode = false;
    },
    edit() {
      this.newListTitle = this.list.title;
      this.editListMode = true;
    },
    saveListTitle() {
      const newTitle = this.newListTitle.trim();
      if (newTitle) {
        const curList = this.list;
        const newList = Object.assign(curList, { title: newTitle });
        this.editList(newList);
        this.editListMode = false;
        this.newListTitle = '';
      }
    },
  },
};
</script>

<style scoped>
</style>
