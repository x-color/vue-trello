<template>
  <div>
    <div v-if="board">
      <!-- PIN: Board info -->
      <v-container>
        <v-row>
          <h1 class="text-center ma-3">{{ board.title }}</h1>
          <div class="text-center my-auto">
            <v-menu offset-y>
              <template v-slot:activator="{ on }">
                <v-btn outlined small text dark v-on="on">
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
        </v-row>
        <v-row style="height: 15px">
          <p class="mx-4">{{ board.text }}</p>
        </v-row>

        <!-- PIN: Modal for editing board -->
        <ModalBoard
          v-model="editedBoard"
          :open="editBoardMode"
          @close="editBoardMode = false"
          @save="saveEditedBoard"
        />
      </v-container>

      <!-- PIN: Lists -->
      <v-container>
        <draggable
          class="row flex-nowrap row--dense justify-start"
          group="lists"
          v-model="lists"
          draggable=".item"
        >
          <v-col v-for="(list, i) in lists" :key="i" cols="auto" class="item">
            <CardList :id="list.id" />
          </v-col>

          <!-- PIN: Button for creating new list -->
          <v-col cols="auto">
            <v-btn
              v-if="!addListMode"
              class="mx-auto"
              width="300"
              height="100"
              color="grey lighten-2"
              @click="addListMode = true"
            >
              <v-icon>mdi-plus</v-icon>
            </v-btn>

            <!-- PIN: Form for new list title -->
            <v-card
              v-if="addListMode"
              class="mx-auto"
              width="300"
              color="grey lighten-4"
            >
              <v-card-title>
                <v-text-field
                  class="headline"
                  v-model="newListTitle"
                  v-if="addListMode"
                  dense
                  full-width
                  label="New List"
                  single-line
                  autofocus
                  @keypress.enter="addNewList()"
                  @blur="addNewList()"
                />
              </v-card-title>
            </v-card>
          </v-col>
        </draggable>
      </v-container>

      <!-- PIN: Confirmation modal for deleting board -->
      <ModalConfirm
        :title="`Delete '${board.title}' ?`"
        text="Can not restore all data in board."
        :open="deleteBoardMode"
        @cancel="deleteBoardMode = false"
        @confirm="del()"
      />
    </div>

    <!-- PIN: No board page -->
    <div v-else>
      <v-container>
        <v-row class="text-center" justify="center">
          <v-col cols="12">
            <h1>This Board does not exist.</h1>
          </v-col>
          <v-col>
            <p>
              Please go to
              <router-link to="/">Home</router-link>
            </p>
          </v-col>
        </v-row>
      </v-container>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import draggable from 'vuedraggable';
import CardList from '@/components/CardList.vue';
import ModalBoard from '@/components/ModalBoard.vue';
import ModalConfirm from '@/components/ModalConfirm.vue';

export default {
  name: 'PageBoard',
  components: {
    CardList,
    ModalBoard,
    ModalConfirm,
    draggable,
  },
  computed: {
    ...mapGetters(['getBoardById', 'getListsByBoardId', 'getListById']),
    board() {
      return this.getBoardById(this.$route.params.id);
    },
    lists: {
      get() {
        return this.getListsByBoardId(this.$route.params.id);
      },
      set(newValue) {
        const newBoard = Object.assign(
          { ...this.board },
          { lists: newValue.map(list => list.id) },
        );
        this.editBoard(newBoard);
      },
    },
  },
  data() {
    return {
      addListMode: false,
      editBoardMode: false,
      deleteBoardMode: false,
      newListTitle: '',
      editedBoard: {},
      menuItems: [
        {
          title: 'edit',
          icon: 'mdi-pencil',
          action: () => {
            this.editedBoard = { ...this.board };
            this.editBoardMode = true;
          },
        },
        {
          title: 'delete',
          icon: 'mdi-trash-can-outline',
          action: () => {
            this.deleteBoardMode = true;
          },
        },
      ],
    };
  },
  methods: {
    ...mapActions(['addList', 'editBoard', 'removeBoard', 'editList']),
    addNewList() {
      const newTitle = this.newListTitle.trim();
      if (newTitle) {
        this.addList({
          boardId: this.board.id,
          title: this.newListTitle.trim(),
        });
      }
      this.newListTitle = '';
      this.addListMode = false;
    },
    saveEditedBoard() {
      this.editBoard(this.editedBoard);
      this.editBoardMode = false;
    },
    del() {
      this.deleteBoardMode = false;
      this.removeBoard(this.board);
    },
  },
};
</script>
