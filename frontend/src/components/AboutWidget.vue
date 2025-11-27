<template>
  <div class="about-page">
    <!-- Hero Section - 个人信息卡片 -->
    <div class="hero-section card">
      <div class="hero-background"></div>
      <div class="hero-content">
        <div class="avatar-wrapper">
          <img :src="profile.avatar_url || '/default-avatar.png'" alt="Avatar" class="hero-avatar" />
          <button v-if="isOwner" @click="openEditModal" class="edit-btn-float" :title="$t('about.edit_profile')">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
            </svg>
          </button>
        </div>
        <h1 class="hero-name">{{ profile.username }}</h1>
        <p class="hero-bio" v-if="profile.bio">{{ profile.bio }}</p>
        <p class="hero-bio-placeholder" v-else-if="isOwner">{{ $t('about.bio_placeholder') }}</p>
      </div>
    </div>

    <!-- Social Links Section -->
    <div class="social-section" v-if="socialLinks.length">
      <h2 class="section-title">
        <span class="title-icon">🔗</span>
        社交链接
      </h2>
      <div class="social-grid">
        <div 
          v-for="(link, index) in socialLinks" 
          :key="index" 
          class="social-card"
          @mouseenter="hoveredSocial = index"
          @mouseleave="hoveredSocial = null"
        >
          <a :href="link.url" target="_blank" class="social-card-link">
            <div class="social-icon" v-html="getIconSvg(link.platform)"></div>
            <span class="social-name">{{ link.name || link.platform }}</span>
          </a>
          <!-- Hover Popup -->
          <div v-if="hoveredSocial === index" class="hover-popup fade-in">
            <div class="popup-url" v-if="link.url">{{ link.url }}</div>
            <img v-if="link.hover_image" :src="link.hover_image" alt="QR Code" />
            <div v-if="!link.hover_image && link.url" class="popup-hint">点击访问</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Sponsor Section -->
    <div class="sponsor-section card" v-if="sponsorLinks.length">
      <div class="sponsor-content">
        <h2 class="section-title">
          <span class="title-icon">❤️</span>
          {{ $t('about.support_my_work') }}
        </h2>
        <p class="sponsor-desc">如果您喜欢我的内容，欢迎赞助支持我的创作</p>
        <div class="sponsor-buttons">
          <div 
            v-for="(link, index) in sponsorLinks" 
            :key="index" 
            class="sponsor-btn-wrapper"
            @mouseenter="hoveredSponsor = index"
            @mouseleave="hoveredSponsor = null"
          >
            <a :href="link.url" target="_blank" class="sponsor-btn">
              <span class="sponsor-platform">{{ link.platform }}</span>
            </a>
            <!-- Hover Popup -->
            <div v-if="hoveredSponsor === index" class="hover-popup fade-in">
              <div class="popup-url" v-if="link.url">{{ link.url }}</div>
              <img v-if="link.hover_image" :src="link.hover_image" alt="Sponsor QR" />
              <div v-if="!link.hover_image && link.url" class="popup-hint">点击赞助</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Friends Section -->
    <div class="friends-section" v-if="friendLinks.length">
      <h2 class="section-title">
        <span class="title-icon">👥</span>
        {{ $t('about.friends') }}
      </h2>
      <div class="friends-grid">
        <div 
          v-for="(friend, index) in friendLinks" 
          :key="index" 
          class="friend-card"
          @mouseenter="hoveredFriend = index"
          @mouseleave="hoveredFriend = null"
        >
          <a :href="friend.url" target="_blank" class="friend-card-link">
            <img :src="friend.avatar || '/default-avatar.png'" class="friend-avatar" />
            <div class="friend-info">
              <h3 class="friend-name">{{ friend.name }}</h3>
              <p class="friend-desc" v-if="friend.description">{{ friend.description }}</p>
            </div>
          </a>
          <!-- Hover Popup -->
          <div v-if="hoveredFriend === index" class="hover-popup friend-hover fade-in">
            <div class="popup-url" v-if="friend.url">{{ friend.url }}</div>
            <img v-if="friend.hover_image" :src="friend.hover_image" alt="Preview" />
            <div v-if="friend.description && !friend.hover_image" class="popup-desc">{{ friend.description }}</div>
            <div v-if="!friend.hover_image && !friend.description && friend.url" class="popup-hint">点击访问</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal" class="modal-overlay" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('about.edit_profile') }}</h3>
          <button type="button" @click="closeEditModal" class="close-btn" title="关闭">✕</button>
        </div>
        <form @submit.prevent="saveProfile">
          <!-- 基本信息 -->
          <div class="form-section">
            <h4 class="section-header">📝 基本信息</h4>
            <div class="form-grid-2">
              <div class="form-group">
                <label>{{ $t('about.avatar_url') }}</label>
                <input v-model="editForm.avatar_url" type="text" placeholder="https://..." />
              </div>
              <div class="form-group full-width">
                <label>{{ $t('about.bio') }}</label>
                <textarea v-model="editForm.bio" rows="3" placeholder="介绍一下自己..."></textarea>
              </div>
            </div>
          </div>

          <!-- 编辑模式切换 -->
          <div class="mode-toggle">
            <button type="button" :class="{ active: editMode === 'visual' }" @click="switchMode('visual')">
              ✏️ {{ $t('about.mode_visual') }}
            </button>
            <button type="button" :class="{ active: editMode === 'code' }" @click="switchMode('code')">
              💻 {{ $t('about.mode_json') }}
            </button>
          </div>

          <!-- 可视化编辑模式 -->
          <div v-if="editMode === 'visual'">
            <!-- 标签页导航 -->
            <div class="tabs-nav">
              <button type="button" :class="{ active: activeTab === 'social' }" @click="activeTab = 'social'">
                🔗 社交链接 ({{ parsedSocialLinks.length }})
              </button>
              <button type="button" :class="{ active: activeTab === 'sponsor' }" @click="activeTab = 'sponsor'">
                ❤️ 赞助链接 ({{ parsedSponsorLinks.length }})
              </button>
              <button type="button" :class="{ active: activeTab === 'friends' }" @click="activeTab = 'friends'">
                👥 友情链接 ({{ parsedFriendLinks.length }})
              </button>
            </div>

            <!-- 社交链接编辑 -->
            <div v-show="activeTab === 'social'" class="tab-content">
              <div class="tab-header">
                <h4>社交链接</h4>
                <button type="button" @click="addSocialLink" class="btn-add-small">+ 添加</button>
              </div>
              <div v-if="parsedSocialLinks.length === 0" class="empty-hint">
                暂无社交链接，点击"添加"按钮创建
              </div>
              <div v-for="(link, index) in parsedSocialLinks" :key="index" class="edit-item">
                <div class="item-header-row" @click="toggleItemExpand('social', index)">
                  <div class="item-number">#{{ index + 1 }}</div>
                  <div class="item-summary">
                    <span class="platform-badge">{{ link.platform }}</span>
                    <span class="item-name">{{ link.name || '未命名' }}</span>
                  </div>
                  <button type="button" class="btn-expand">
                    {{ isItemExpanded('social', index) ? '▼ 收起' : '▶ 展开' }}
                  </button>
                </div>
                <div v-show="isItemExpanded('social', index)" class="item-details">
                  <div class="form-grid-2">
                    <div class="form-group">
                      <label>平台</label>
                      <select v-model="link.platform">
                        <option value="github">GitHub</option>
                        <option value="twitter">Twitter</option>
                        <option value="email">Email</option>
                        <option value="linkedin">LinkedIn</option>
                        <option value="youtube">YouTube</option>
                        <option value="instagram">Instagram</option>
                        <option value="facebook">Facebook</option>
                        <option value="telegram">Telegram</option>
                        <option value="tiktok">TikTok</option>
                        <option value="qq">QQ</option>
                        <option value="wechat">WeChat</option>
                        <option value="bilibili">Bilibili</option>
                        <option value="douyin">Douyin</option>
                      </select>
                    </div>
                    <div class="form-group">
                      <label>显示名称（可选）</label>
                      <input v-model="link.name" placeholder="例如：我的 GitHub" />
                    </div>
                    <div class="form-group full-width">
                      <label>链接地址</label>
                      <input v-model="link.url" placeholder="https://..." />
                    </div>
                    <div class="form-group full-width">
                      <label>悬停图片（可选，如二维码）</label>
                      <input v-model="link.hover_image" placeholder="https://..." />
                    </div>
                  </div>
                  <button type="button" @click="removeSocialLink(index)" class="btn-remove-item">🗑️ 删除</button>
                </div>
              </div>
            </div>

            <!-- 赞助链接编辑 -->
            <div v-show="activeTab === 'sponsor'" class="tab-content">
              <div class="tab-header">
                <h4>赞助链接</h4>
                <button type="button" @click="addSponsorLink" class="btn-add-small">+ 添加</button>
              </div>
              <div v-if="parsedSponsorLinks.length === 0" class="empty-hint">
                暂无赞助链接，点击"添加"按钮创建
              </div>
              <div v-for="(link, index) in parsedSponsorLinks" :key="index" class="edit-item">
                <div class="item-header-row" @click="toggleItemExpand('sponsor', index)">
                  <div class="item-number">#{{ index + 1 }}</div>
                  <div class="item-summary">
                    <span class="platform-badge">{{ link.platform }}</span>
                  </div>
                  <button type="button" class="btn-expand">
                    {{ isItemExpanded('sponsor', index) ? '▼ 收起' : '▶ 展开' }}
                  </button>
                </div>
                <div v-show="isItemExpanded('sponsor', index)" class="item-details">
                  <div class="form-grid-2">
                    <div class="form-group">
                      <label>平台名称</label>
                      <input v-model="link.platform" placeholder="例如：Patreon" />
                    </div>
                    <div class="form-group">
                      <label>链接地址</label>
                      <input v-model="link.url" placeholder="https://..." />
                    </div>
                  </div>
                  <button type="button" @click="removeSponsorLink(index)" class="btn-remove-item">🗑️ 删除</button>
                </div>
              </div>
            </div>

            <!-- 友情链接编辑 -->
            <div v-show="activeTab === 'friends'" class="tab-content">
              <div class="tab-header">
                <h4>友情链接</h4>
                <button type="button" @click="addFriendLink" class="btn-add-small">+ 添加</button>
              </div>
              <div v-if="parsedFriendLinks.length === 0" class="empty-hint">
                暂无友情链接，点击"添加"按钮创建
              </div>
              <div v-for="(link, index) in parsedFriendLinks" :key="index" class="edit-item">
                <div class="item-header-row" @click="toggleItemExpand('friends', index)">
                  <div class="item-number">#{{ index + 1 }}</div>
                  <div class="item-summary">
                    <span class="item-name">{{ link.name || '未命名' }}</span>
                  </div>
                  <button type="button" class="btn-expand">
                    {{ isItemExpanded('friends', index) ? '▼ 收起' : '▶ 展开' }}
                  </button>
                </div>
                <div v-show="isItemExpanded('friends', index)" class="item-details">
                  <div class="form-grid-2">
                    <div class="form-group">
                      <label>名称</label>
                      <input v-model="link.name" placeholder="好友名称" />
                    </div>
                    <div class="form-group">
                      <label>链接地址</label>
                      <input v-model="link.url" placeholder="https://..." />
                    </div>
                    <div class="form-group full-width">
                      <label>头像地址</label>
                      <input v-model="link.avatar" placeholder="https://..." />
                    </div>
                    <div class="form-group full-width">
                      <label>描述（可选）</label>
                      <input v-model="link.description" placeholder="简短介绍..." />
                    </div>
                    <div class="form-group full-width">
                      <label>悬停图片（可选）</label>
                      <input v-model="link.hover_image" placeholder="https://..." />
                    </div>
                  </div>
                  <button type="button" @click="removeFriendLink(index)" class="btn-remove-item">🗑️ 删除</button>
                </div>
              </div>
            </div>
          </div>

          <!-- JSON 代码模式 -->
          <div v-else class="code-mode">
            <div class="form-group">
              <label>{{ $t('about.social_links') }}</label>
              <div class="help-text" v-html="$t('about.social_help')"></div>
              <textarea v-model="editForm.social_links" rows="8" class="code-input"></textarea>
            </div>
            <div class="form-group">
              <label>{{ $t('about.sponsor_links') }}</label>
              <div class="help-text" v-html="$t('about.sponsor_help')"></div>
              <textarea v-model="editForm.sponsor_links" rows="5" class="code-input"></textarea>
            </div>
            <div class="form-group">
              <label>{{ $t('about.friend_links') }}</label>
              <div class="help-text" v-html="$t('about.friend_help')"></div>
              <textarea v-model="editForm.friend_links" rows="8" class="code-input"></textarea>
            </div>
          </div>

          <div class="modal-actions">
            <button type="button" @click="closeEditModal" class="btn btn-ghost">{{ $t('common.cancel') }}</button>
            <button type="submit" class="btn btn-primary">{{ $t('common.save') }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Sponsor Modal -->
    <div v-if="showSponsorModal" class="modal-overlay" @click.self="showSponsorModal = false">
      <div class="modal-content card text-center">
        <h3>Support My Work</h3>
        <div class="sponsor-list">
          <a v-for="(link, index) in sponsorLinks" :key="index" :href="link.url" target="_blank" class="btn btn-outline w-full mb-2">
            {{ link.platform }}
          </a>
        </div>
        <button @click="showSponsorModal = false" class="btn btn-ghost mt-4">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import api from '../api/axios';
import { useToast } from '../composables/useToast';

const authStore = useAuthStore();
const toast = useToast();
const profile = ref({});
const showEditModal = ref(false);
const showSponsorModal = ref(false);
const hoveredSocial = ref(null);
const hoveredSponsor = ref(null);
const hoveredFriend = ref(null);

const editForm = ref({
  avatar_url: '',
  bio: '',
  social_links: '[]',
  sponsor_links: '[]',
  friend_links: '[]'
});

const editMode = ref('visual'); // 'visual' or 'code'
const activeTab = ref('social'); // 'social', 'sponsor', 'friends'
const expandedItems = ref({
  social: [],
  sponsor: [],
  friends: []
});
const parsedSocialLinks = ref([]);
const parsedSponsorLinks = ref([]);
const parsedFriendLinks = ref([]);

const isOwner = computed(() => {
  return authStore.isAuthenticated;
});

const socialLinks = computed(() => {
  try {
    const parsed = JSON.parse(profile.value.social_links || '[]');
    console.log('Parsed social links:', parsed);
    return parsed;
  } catch (e) {
    console.error('Error parsing social links:', e);
    return [];
  }
});

const sponsorLinks = computed(() => {
  try {
    const parsed = JSON.parse(profile.value.sponsor_links || '[]');
    console.log('Parsed sponsor links:', parsed);
    return parsed;
  } catch (e) {
    console.error('Error parsing sponsor links:', e);
    return [];
  }
});

const friendLinks = computed(() => {
  try {
    const parsed = JSON.parse(profile.value.friend_links || '[]');
    console.log('Parsed friend links:', parsed);
    return parsed;
  } catch (e) {
    console.error('Error parsing friend links:', e);
    return [];
  }
});

const getIconSvg = (platform) => {
  const p = platform.toLowerCase();
  const svgs = {
    github: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>',
    twitter: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"/></svg>',
    email: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"/></svg>',
    youtube: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/></svg>',
    instagram: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/></svg>',
    facebook: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/></svg>',
    telegram: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 11.944 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z"/></svg>',
    tiktok: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z"/></svg>',
    linkedin: '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.79-1.75-1.764s.784-1.764 1.75-1.764 1.75.79 1.75 1.764-.783 1.764-1.75 1.764zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z"/></svg>',
    // Simple text-based fallbacks or SVGs for Chinese platforms if not available standardly, 
    // using generic shapes or letters for now to ensure they render something.
    // Ideally we'd have real SVGs for these. I'll use placeholders with text inside SVG for clarity.
    qq: '<svg viewBox="0 0 24 24" fill="currentColor"><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="14" font-weight="bold">QQ</text></svg>',
    wechat: '<svg viewBox="0 0 24 24" fill="currentColor"><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="10" font-weight="bold">WX</text></svg>',
    bilibili: '<svg viewBox="0 0 24 24" fill="currentColor"><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="10" font-weight="bold">Bili</text></svg>',
    douyin: '<svg viewBox="0 0 24 24" fill="currentColor"><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="10" font-weight="bold">DY</text></svg>',
  };
  return svgs[p] || `<svg viewBox="0 0 24 24" fill="currentColor"><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="12" font-weight="bold">${p.substring(0, 2).toUpperCase()}</text></svg>`;
};

const fetchProfile = async () => {
  try {
    // If user is logged in, get their profile
    if (authStore.isAuthenticated) {
      const response = await api.get('/user/profile');
      console.log('Fetched user profile:', response.data);
      profile.value = response.data;
    } else {
      // If not logged in, get public author profile
      const response = await api.get('/author-profile');
      console.log('Fetched author profile:', response.data);
      profile.value = response.data;
    }
  } catch (err) {
    console.error('Failed to fetch profile:', err);
    // Fallback to public profile if user profile fails
    try {
      const response = await api.get('/author-profile');
      profile.value = response.data;
    } catch (fallbackErr) {
      console.error('Failed to fetch fallback profile:', fallbackErr);
    }
  }
};

const openEditModal = () => {
  editForm.value = {
    avatar_url: profile.value.avatar_url || '',
    bio: profile.value.bio || '',
    social_links: profile.value.social_links || '[]',
    sponsor_links: profile.value.sponsor_links || '[]',
    friend_links: profile.value.friend_links || '[]'
  };
  
  // Initialize parsed arrays for visual editor
  try {
    parsedSocialLinks.value = JSON.parse(editForm.value.social_links);
  } catch { parsedSocialLinks.value = []; }
  
  try {
    parsedSponsorLinks.value = JSON.parse(editForm.value.sponsor_links);
  } catch { parsedSponsorLinks.value = []; }
  
  try {
    parsedFriendLinks.value = JSON.parse(editForm.value.friend_links);
  } catch { parsedFriendLinks.value = []; }

  showEditModal.value = true;
};

const closeEditModal = () => {
  showEditModal.value = false;
};

const switchMode = (mode) => {
  if (mode === 'visual') {
    // Sync JSON strings to arrays
    try {
      parsedSocialLinks.value = JSON.parse(editForm.value.social_links);
      parsedSponsorLinks.value = JSON.parse(editForm.value.sponsor_links);
      parsedFriendLinks.value = JSON.parse(editForm.value.friend_links);
      editMode.value = 'visual';
    } catch (e) {
      toast.error(t('about.json_error'));
    }
  } else {
    // Sync arrays to JSON strings
    editForm.value.social_links = JSON.stringify(parsedSocialLinks.value, null, 2);
    editForm.value.sponsor_links = JSON.stringify(parsedSponsorLinks.value, null, 2);
    editForm.value.friend_links = JSON.stringify(parsedFriendLinks.value, null, 2);
    editMode.value = 'code';
  }
};

const addSocialLink = () => {
  parsedSocialLinks.value.push({ platform: 'github', name: '', url: '', hover_image: '' });
};

const removeSocialLink = (index) => {
  parsedSocialLinks.value.splice(index, 1);
};

const addSponsorLink = () => {
  parsedSponsorLinks.value.push({ platform: 'patreon', url: '' });
};

const removeSponsorLink = (index) => {
  parsedSponsorLinks.value.splice(index, 1);
};

const addFriendLink = () => {
  parsedFriendLinks.value.push({ name: '', url: '', avatar: '', description: '', hover_image: '' });
};

const removeFriendLink = (index) => {
  parsedFriendLinks.value.splice(index, 1);
};

const toggleItemExpand = (type, index) => {
  const idx = expandedItems.value[type].indexOf(index);
  if (idx > -1) {
    expandedItems.value[type].splice(idx, 1);
  } else {
    expandedItems.value[type].push(index);
  }
};

const isItemExpanded = (type, index) => {
  return expandedItems.value[type].includes(index);
};

const saveProfile = async () => {
  try {
    // If in visual mode, sync back to strings first
    if (editMode.value === 'visual') {
      editForm.value.social_links = JSON.stringify(parsedSocialLinks.value);
      editForm.value.sponsor_links = JSON.stringify(parsedSponsorLinks.value);
      editForm.value.friend_links = JSON.stringify(parsedFriendLinks.value);
    } else {
      // Validate JSON if in code mode
      JSON.parse(editForm.value.social_links);
      JSON.parse(editForm.value.sponsor_links);
      JSON.parse(editForm.value.friend_links);
    }

    console.log('Saving profile:', editForm.value);
    
    const response = await api.put('/user/profile', editForm.value);
    console.log('Save response:', response.data);
    
    // Update local profile with response data
    profile.value = response.data;
    
    // Close modal
    showEditModal.value = false;
    
    // Show success message
    toast.success('Profile updated!');
    
    // Refetch to ensure we have the latest data
    await fetchProfile();
  } catch (err) {
    console.error('Save error:', err);
    if (err instanceof SyntaxError) {
      toast.error('Invalid JSON format');
    } else {
      toast.error('Failed to update profile');
    }
  }
};

onMounted(() => {
  fetchProfile();
});
</script>

<style scoped>
/* About Page Layout */
.about-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

/* Hero Section */
.hero-section {
  position: relative;
  overflow: hidden;
  margin-bottom: 3rem;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(0, 243, 255, 0.1) 0%, rgba(138, 43, 226, 0.1) 100%);
}

.hero-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: linear-gradient(135deg, var(--primary-color) 0%, #8a2be2 100%);
  opacity: 0.3;
  z-index: 0;
}

.hero-content {
  position: relative;
  z-index: 1;
  padding: 3rem 2rem 2rem;
  text-align: center;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 1.5rem;
}

.hero-avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  border: 5px solid var(--primary-color);
  box-shadow: 0 0 30px rgba(0, 243, 255, 0.5),
              0 10px 40px rgba(0, 0, 0, 0.3);
  transition: transform 0.3s;
}

.hero-avatar:hover {
  transform: scale(1.05);
}

.edit-btn-float {
  position: absolute;
  bottom: 5px;
  right: 5px;
  background: var(--secondary-color);
  color: #000;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transition: all 0.3s;
}

.edit-btn-float svg {
  width: 20px;
  height: 20px;
}

.edit-btn-float:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}

.hero-name {
  font-size: 2.5rem;
  color: var(--primary-color);
  font-family: var(--font-heading);
  text-shadow: var(--neon-text-shadow);
  margin: 0 0 1rem 0;
}

.hero-bio {
  font-size: 1.1rem;
  color: var(--text-secondary);
  line-height: 1.8;
  max-width: 600px;
  margin: 0 auto;
}

.hero-bio-placeholder {
  font-size: 1rem;
  color: var(--text-secondary);
  font-style: italic;
  opacity: 0.6;
}

/* Section Title */
.section-title {
  font-size: 1.8rem;
  color: var(--primary-color);
  font-family: var(--font-heading);
  margin-bottom: 2rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.title-icon {
  font-size: 2rem;
}

/* Social Section */
.social-section {
  margin-bottom: 3rem;
}

.social-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1.5rem;
}

.social-card {
  position: relative;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  padding: 1.5rem;
  text-align: center;
  transition: all 0.3s;
}

.social-card:hover {
  background: rgba(0, 243, 255, 0.1);
  border-color: var(--primary-color);
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 243, 255, 0.2);
}

.social-card-link {
  text-decoration: none;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}

.social-icon {
  width: 48px;
  height: 48px;
  color: var(--primary-color);
  transition: transform 0.3s;
}

.social-icon svg {
  width: 100%;
  height: 100%;
}

.social-card:hover .social-icon {
  transform: scale(1.2);
}

.social-name {
  font-size: 0.9rem;
  color: var(--text-secondary);
  font-weight: 500;
}

/* Sponsor Section */
.sponsor-section {
  margin-bottom: 3rem;
  background: linear-gradient(135deg, rgba(255, 0, 128, 0.1) 0%, rgba(255, 128, 0, 0.1) 100%);
  border-radius: 16px;
  overflow: hidden;
}

.sponsor-content {
  padding: 2.5rem 2rem;
  text-align: center;
}

.sponsor-desc {
  font-size: 1rem;
  color: var(--text-secondary);
  margin-bottom: 2rem;
}

.sponsor-buttons {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.sponsor-btn {
  background: linear-gradient(135deg, #ff0080 0%, #ff8000 100%);
  color: white;
  text-decoration: none;
  padding: 12px 30px;
  border-radius: 25px;
  font-weight: bold;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(255, 0, 128, 0.3);
}

.sponsor-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(255, 0, 128, 0.5);
}

/* Friends Section */
.friends-section {
  margin-bottom: 3rem;
}

.friends-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

.friend-card {
  position: relative;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
}

.friend-card:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: var(--secondary-color);
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(138, 43, 226, 0.2);
}

.friend-card-link {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  text-decoration: none;
}

.friend-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--secondary-color);
  flex-shrink: 0;
  transition: transform 0.3s;
}

.friend-card:hover .friend-avatar {
  transform: scale(1.1);
}

.friend-info {
  flex: 1;
  min-width: 0;
}

.friend-name {
  font-size: 1.1rem;
  color: var(--primary-color);
  margin: 0 0 0.5rem 0;
  font-weight: 600;
}

.friend-desc {
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Hover Popup */
.hover-popup {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  padding: 8px;
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.3);
  z-index: 100;
  margin-bottom: 15px;
  width: 150px;
  animation: fadeIn 0.3s;
}

.hover-popup img {
  width: 100%;
  display: block;
  border-radius: 4px;
}

.hover-popup::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  margin-left: -8px;
  border-width: 8px;
  border-style: solid;
  border-color: #fff transparent transparent transparent;
}

.friend-hover {
  width: 200px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.fade-in {
  animation: fadeIn 0.3s;
}

/* Modal Styles - Full Screen */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.95);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
  padding: 0;
}

.modal-content {
  width: 100%;
  height: 100%;
  max-width: none;
  max-height: none;
  overflow-y: auto;
  background: #0a0a12;
  border: none;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.modal-content h3 {
  margin-bottom: 2rem;
  font-size: 2rem;
  text-align: center;
  background: #0a0a12;
  padding: 1rem 0;
  margin: 0;
}

.modal-header {
  background: #0a0a12;
  padding: 2rem 3rem;
  border-bottom: 1px solid var(--primary-color);
  margin-bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.close-btn {
  position: absolute;
  right: 3rem;
  top: 2rem;
  background: transparent;
  border: 2px solid var(--primary-color);
  color: var(--primary-color);
  width: 50px;
  height: 50px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  line-height: 1;
}

.close-btn:hover {
  background: var(--primary-color);
  color: #000;
  transform: rotate(90deg);
}

.modal-content form {
  width: 100%;
  padding: 2rem 5%;
  flex: 1;
  max-width: 1400px;
  margin: 0 auto;
}

.form-group {
  margin-bottom: 2rem;
  text-align: left;
}

.form-group label {
  display: block;
  color: var(--primary-color);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.help-text {
  font-size: 0.75rem;
  color: var(--text-secondary);
  margin-bottom: 0.75rem;
  line-height: 1.6;
  padding: 8px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 4px;
}

.help-text code {
  background: rgba(255,255,255,0.1);
  padding: 2px 6px;
  border-radius: 2px;
  font-family: monospace;
  font-size: 0.7rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
  padding: 0.5rem;
  border-radius: 4px;
}

.code-input {
  font-family: monospace;
  font-size: 0.85rem;
  line-height: 1.5;
  min-height: 120px;
}

.mode-toggle {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
  gap: 15px;
}

.mode-toggle button {
  background: transparent;
  border: 1px solid var(--primary-color);
  color: var(--text-secondary);
  padding: 8px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s;
}

.mode-toggle button:hover {
  background: rgba(0, 243, 255, 0.1);
}

.mode-toggle button.active {
  background: var(--primary-color);
  color: #000;
  font-weight: bold;
}

.visual-editor {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-top: 15px;
}

.visual-item {
  background: rgba(255, 255, 255, 0.05);
  padding: 20px;
  border-radius: 8px;
  border: 1px solid var(--surface-border);
  position: relative;
  transition: all 0.3s;
}

.visual-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--primary-color);
}

.item-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
  margin-bottom: 12px;
}

.item-row:last-child {
  margin-bottom: 0;
}

.input-sm {
  min-width: 0;
}

.input-full {
  width: 100%;
  grid-column: 1 / -1;
}

.form-group input,
.form-group textarea {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

.btn-remove {
  background: #ff4d4f;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 0.8rem;
  cursor: pointer;
  margin-top: 8px;
  transition: background 0.3s;
}

.btn-remove:hover {
  background: #ff7875;
}

.btn-add {
  background: transparent;
  border: 1px dashed var(--primary-color);
  color: var(--primary-color);
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  font-size: 0.9rem;
  transition: all 0.3s;
}

.btn-add:hover {
  background: rgba(0, 243, 255, 0.1);
  border-style: solid;
}

.modal-actions {
  display: flex;
  justify-content: center;
  gap: 2rem;
  margin-top: 3rem;
  padding: 2rem 0;
}

.modal-actions button {
  min-width: 150px;
  padding: 12px 30px;
  font-size: 1rem;
}

.w-full { width: 100%; }
.mb-2 { margin-bottom: 0.5rem; }
.mt-4 { margin-top: 1rem; }

/* Form Sections */
.form-section {
  margin-bottom: 2.5rem;
}

.section-header {
  font-size: 1.2rem;
  color: var(--primary-color);
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid var(--surface-border);
}

.form-grid-2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

/* Tabs Navigation */
.tabs-nav {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  border-bottom: 2px solid var(--surface-border);
}

.tabs-nav button {
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  color: var(--text-secondary);
  padding: 1rem 1.5rem;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
  margin-bottom: -2px;
}

.tabs-nav button:hover {
  color: var(--primary-color);
  background: rgba(0, 243, 255, 0.05);
}

.tabs-nav button.active {
  color: var(--primary-color);
  border-bottom-color: var(--primary-color);
  font-weight: bold;
}

/* Tab Content */
.tab-content {
  animation: fadeIn 0.3s;
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.tab-header h4 {
  margin: 0;
  color: var(--secondary-color);
  font-size: 1.1rem;
}

.btn-add-small {
  background: var(--primary-color);
  color: #000;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s;
}

.btn-add-small:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 243, 255, 0.4);
}

/* Edit Items */
.edit-item {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  padding: 0;
  margin-bottom: 1.5rem;
  transition: all 0.3s;
  overflow: hidden;
}

.edit-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: var(--primary-color);
}

.item-header-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.5rem;
  cursor: pointer;
  transition: background 0.3s;
  position: relative;
}

.item-header-row:hover {
  background: rgba(0, 243, 255, 0.05);
}

.item-summary {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.platform-badge {
  background: var(--primary-color);
  color: #000;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: bold;
  text-transform: uppercase;
}

.item-name {
  color: var(--text-primary);
  font-size: 1rem;
}

.btn-expand {
  background: transparent;
  border: 1px solid var(--surface-border);
  color: var(--text-secondary);
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.3s;
  white-space: nowrap;
}

.btn-expand:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.item-details {
  padding: 0 1.5rem 1.5rem 1.5rem;
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    max-height: 0;
  }
  to {
    opacity: 1;
    max-height: 1000px;
  }
}

.item-number {
  background: var(--primary-color);
  color: #000;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: bold;
}

.btn-remove-item {
  background: #ff4d4f;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  margin-top: 1rem;
  transition: all 0.3s;
}

.btn-remove-item:hover {
  background: #ff7875;
  transform: translateX(5px);
}

/* Empty Hint */
.empty-hint {
  text-align: center;
  padding: 3rem 2rem;
  color: var(--text-secondary);
  font-style: italic;
  background: rgba(255, 255, 255, 0.02);
  border: 2px dashed var(--surface-border);
  border-radius: 8px;
}

/* Code Mode */
.code-mode {
  animation: fadeIn 0.3s;
}

/* Form Inputs Enhancement */
.form-group select,
.form-group input,
.form-group textarea {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--surface-border);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 6px;
  font-size: 0.95rem;
  transition: all 0.3s;
}

.form-group select:focus,
.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(0, 243, 255, 0.1);
}

.form-group label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
}

/* Fix line-clamp compatibility */
.friend-desc {
  line-clamp: 2;
  -webkit-line-clamp: 2;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Hover Popup Styles */
.hover-popup {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-top: 10px;
  background: rgba(10, 10, 18, 0.98);
  border: 1px solid var(--primary-color);
  border-radius: 8px;
  padding: 1rem;
  min-width: 200px;
  max-width: 300px;
  z-index: 1000;
  box-shadow: 0 8px 32px rgba(0, 243, 255, 0.3);
  pointer-events: none;
}

.hover-popup::before {
  content: '';
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 8px solid transparent;
  border-bottom-color: var(--primary-color);
}

.popup-url {
  color: var(--primary-color);
  font-size: 0.85rem;
  margin-bottom: 0.75rem;
  word-break: break-all;
  line-height: 1.4;
  font-family: monospace;
}

.hover-popup img {
  width: 100%;
  max-width: 200px;
  height: auto;
  border-radius: 4px;
  display: block;
  margin: 0 auto;
}

.popup-hint {
  color: var(--text-secondary);
  font-size: 0.9rem;
  text-align: center;
  font-style: italic;
  margin-top: 0.5rem;
}

.popup-desc {
  color: var(--text-secondary);
  font-size: 0.9rem;
  line-height: 1.5;
  margin-top: 0.5rem;
}

.friend-hover {
  min-width: 250px;
  max-width: 350px;
}

/* Sponsor Button Wrapper */
.sponsor-btn-wrapper {
  position: relative;
  display: inline-block;
}

.sponsor-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}
</style>
