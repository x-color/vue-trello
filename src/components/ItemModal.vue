<template>
  <v-dialog v-model="open" persistent max-width="600px" @click:outside="close">
    <v-card>
      <v-container>
        <v-row>
          <!-- PIN: Title -->
          <v-col cols="12">
            <v-card-title
              v-if="!editTitleMode"
              class="headline pt-1 pb-0"
              :class="{ 'font-weight-light': !value.title }"
              @click="editTitleMode = true"
            >{{ value.title | replaceToHintTitle }}</v-card-title>
            <v-card-title v-else class="py-0">
              <v-text-field
                class="headline pt-1 pb-0"
                dense
                hide-details
                v-model="value.title"
                autofocus
                @blur="editTitleMode = false"
              ></v-text-field>
            </v-card-title>
          </v-col>

          <!-- PIN: Text -->
          <v-col cols="12" class="pt-0">
            <v-card-text
              v-if="!editTextMode"
              class="pt-0 body-2"
              :class="{ 'font-weight-light': !value.text }"
              @click="editTextMode = true"
            >{{ value.text | replaceToHintText }}</v-card-text>
            <v-card-text v-else class="pt-0 pb-2">
              <v-textarea
                class="body-2"
                v-model="value.text"
                autofocus
                placeholder="Add desription..."
                auto-grow
                dense
                hide-details
                @blur="editTextMode = false"
              ></v-textarea>
            </v-card-text>
          </v-col>

          <!-- PIN: Tags -->
          <v-col cols="12" class="px-5">
            <!--
              NOTE: This `span` tag lines tags up. (`tag1`, `tag2`, `+`)
                    If remove it, tags line up `+`, `tag1`, `tag2`.
            -->
            <span>
              <BaseTag
                :tag="tag"
                v-for="(tag, index) in tags"
                class="ma-1"
                :key="index"
              />
            </span>

            <!-- PIN: Tags Menu -->
            <v-menu
              v-model="isOpenedMenu"
              :close-on-content-click="false"
              offset-y
            >
              <template v-slot:activator="{ on }">
                <v-chip class="ma-2" color="gray" small label v-on="on">
                  <v-icon small>mdi-plus</v-icon>
                </v-chip>
              </template>
              <v-card>
                <v-list dense subheader max-width="400">
                  <v-subheader>Tags</v-subheader>
                  <v-list-item
                    v-for="(tag, index) in $store.state.tag.tags"
                    :key="index"
                    class="px-2"
                  >
                    <v-chip
                      :color="tag.color"
                      label
                      text-color="white"
                      style="width: 100%"
                      class="pl-0"
                      @click="toggleTag(tag)"
                    >
                      <v-container fluid>
                        <v-row justify="start" no-gutters>
                          <v-col cols="auto">
                            <v-icon v-if="isActive(tag)" small left>mdi-check</v-icon>
                            <v-icon
                              v-else
                              :color="tag.color"
                              small
                              left
                            >mdi-check</v-icon>
                            {{ tag.title }}
                          </v-col>
                        </v-row>
                      </v-container>
                    </v-chip>
                  </v-list-item>
                </v-list>
              </v-card>
            </v-menu>
          </v-col>
        </v-row>

        <v-row justify="end">
          <v-col cols="auto">
            <v-btn
              v-if="isHoverDeleteIcon"
              text
              small
              icon
              color="red"
              @mouseleave="isHoverDeleteIcon = false"
              @click="deleteItemMode = true"
            >
              <v-icon>mdi-delete</v-icon>
            </v-btn>
            <v-btn
              v-else
              text
              small
              icon
              color="red"
              @mouseover="isHoverDeleteIcon = true"
            >
              <v-icon>mdi-delete-outline</v-icon>
            </v-btn>
          </v-col>
          <v-col cols="auto" class="mr-3">
            <v-btn
              v-if="isHoverCheckIcon"
              text
              small
              icon
              color="green"
              @mouseleave="isHoverCheckIcon = false"
              @click="close"
            >
              <v-icon>mdi-check-bold</v-icon>
            </v-btn>
            <v-btn
              v-else
              text
              small
              icon
              color="green"
              @mouseover="isHoverCheckIcon = true"
            >
              <v-icon>mdi-check</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!-- PIN: Confirmation modal for deleting item -->
    <ConfirmModal
      :title="`Delete '${value.title}' ?`"
      text="Can not restore this item."
      :open="deleteItemMode"
      @cancel="deleteItemMode = false"
      @confirm="del()"
    />
  </v-dialog>
</template>

<script>
import Vue from 'vue';
import { mapGetters } from 'vuex';
import ConfirmModal from '@/components/ConfirmModal.vue';
import BaseTag from '@/components/BaseTag.vue';

export default {
  name: 'item-modal',
  props: {
    value: Object,
    open: Boolean,
  },
  components: {
    ConfirmModal,
    BaseTag,
  },
  computed: {
    ...mapGetters(['getTagById']),
    tags() {
      return this.value.tags.map(tagId => this.getTagById(tagId));
    },
  },
  data() {
    return {
      editTitleMode: false,
      editTextMode: false,
      deleteItemMode: false,
      isOpenedMenu: false,
      isHoverDeleteIcon: false,
      isHoverCheckIcon: false,
      activeList: {},
    };
  },
  methods: {
    addTag(id) {
      this.value.tags.push(id);
    },
    removeTag(id) {
      this.value.tags = this.value.tags.filter(tagId => tagId !== id);
    },
    isActive({ id }) {
      if (id in this.activeList) {
        return this.activeList[id];
      }
      if (this.value.tags.find(tagId => tagId === id)) {
        Vue.set(this.activeList, id, true);
        return true;
      }
      return false;
    },
    toggleTag({ id }) {
      if (id in this.activeList) {
        Vue.set(this.activeList, id, !this.activeList[id]);
      } else {
        Vue.set(this.activeList, id, true);
      }
      if (this.activeList[id]) {
        this.addTag(id);
      } else {
        this.removeTag(id);
      }
    },
    close() {
      if (!this.isOpenedMenu) {
        this.$emit('close');
      }
    },
    del() {
      this.deleteItemMode = false;
      this.$emit('delete');
    },
  },
  filters: {
    replaceToHintText(text) {
      if (!text) {
        return 'Add description...';
      }
      return text;
    },
    replaceToHintTitle(text) {
      if (!text) {
        return 'No title...';
      }
      return text;
    },
  },
};
</script>
