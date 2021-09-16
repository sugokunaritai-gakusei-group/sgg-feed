"use strict";

class DOMEV {
    static getDOM() {
        this.select = document.getElementById("select");
        this.urlInput = document.getElementById("header_url");
        this.articlesContainer = document.getElementById("articles_container");
    }
    static changeURL() {
        this.urlInput.value = this.select.value;
    }
    static async getArticles() {
        const res = await fetch("/api");
        const data = await res.json();
        const items = data.items;
        items.sort((a, b) => {
            if (a.date_published > b.date_published) return -1;
            return 1;
        });
        return items;
    }
    static generateHTML(item) {
        const date = item.date_published.split("T")[0].replace(/-/g, "/");
        let site = item.url.match(/^https?:\/{2,}(.*?)(?:\/|\?|#|$)/)[1];
        const structure = `
        <a href="${item.url}" class="links">
            <div class="articles-container-item">
                <div class="articles-container-item-title">${item.title}</div>
                <div class="articles-container-item-info">
                    <div class="articles-container-item-info-date">${date}</div>
                    <div class="articles-container-item-info-site">${site}</div>
                </div>
            </div>
        </a>
    `;
        this.articlesContainer.insertAdjacentHTML("beforeend", structure);
    }
}
window.addEventListener("DOMContentLoaded", async () => {
    DOMEV.getDOM();
    const items = await DOMEV.getArticles();
    for (let i = 0; i < 50; i++) {
        DOMEV.generateHTML(items[i]);
    }
})
