<template>
  <div>
    <div v-if="board">
      <!-- Board info -->
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
        <BoardModal
          v-model="editedBoard"
          :open="editBoardMode"
          @close="editBoardMode = false"
          @save="saveEditedBoard"
        ></BoardModal>
      </v-container>

      <!-- Lists -->
      <v-container>
        <v-row class="flex-nowrap" dense justify="start">
          <v-col v-for="(list, i) in lists" :key="i" cols="auto">
            <List :id="list.id" />
          </v-col>

          <!-- Add new list -->
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

            <!-- New list form -->
            <v-card v-if="addListMode" class="mx-auto" width="300" color="grey lighten-4">
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
                ></v-text-field>
              </v-card-title>
            </v-card>
          </v-col>
        </v-row>
      </v-container>

      <ConfirmModal
        :title="`Delete '${board.title}' ?`"
        text="Can not restore all data in board."
        :open="deleteBoardMode"
        @cancel="deleteBoardMode = false"
        @confirm="del()"
      />
    </div>
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
import List from '@/components/List.vue';
import BoardModal from '@/components/BoardModal.vue';
import ConfirmModal from '@/components/ConfirmModal.vue';

export default {
  name: 'board-page',
  components: {
    List,
    BoardModal,
    ConfirmModal,
  },
  computed: {
    board() {
      return this.getBoardById()(this.$route.params.id);
    },
    lists() {
      return this.getListsByBoardId()(this.$route.params.id);
    },
  },
  methods: {
    ...mapGetters(['getBoardById', 'getListsByBoardId']),
    ...mapActions(['addList', 'editBoard', 'removeBoard']),
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
      this.editBoard({
        id: this.board.id,
        userId: this.board.userId,
        color: this.editedBoard.color,
        text: this.editedBoard.text,
        title: this.editedBoard.title,
      });
      this.editBoardMode = false;
    },
    del() {
      this.deleteBoardMode = false;
      this.removeBoard(this.board);
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
};
</script>
