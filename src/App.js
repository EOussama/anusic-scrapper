import React, { Component } from 'react'
import Axios from 'axios'

import './App.css';
import 'bootstrap/dist/css/bootstrap.css';

import AnimeInfo from './components/AnimeInfo/AnimeInfo';

export default class App extends Component {

  //#region Properties

  state = {
    animeList: [],
    list: [],
    anime: {},
    infoShown: false
  }

  //#endregion

  //#region Lifecycle

  componentDidMount() {
    Axios.get('https://anusic-api.herokuapp.com/api/v1/anime')
      .then(e => {
        this.setState({ animeList: e.data, list: e.data });
      });
  }

  render() {
    return (
      <React.Fragment>
        <nav className="navbar sticky-top navbar-light bg-light px-5">
          <a href="/#" className="navbar-brand">Anusic React</a>
          <div className="form-inline">
            <input
              className="form-control mr-sm-2"
              type="search"
              placeholder="Search"
              aria-label="Search"
              onChange={(e) => {
                if (e.target.value.length > 0) {
                  this.setState({
                    list: this.state.animeList.filter(anime =>
                      anime.name.toLowerCase()
                        .concat((anime.altNames || []).join(' ').toLowerCase())
                        .concat((anime.year || 0).toString())
                        .includes(e.target.value.toLowerCase().trim())
                    )
                  });
                } else {
                  this.setState({ list: this.state.animeList });
                }
              }} />
          </div>
        </nav>

        <main className="p-5">
          <div
            className="alert alert-dark"
            role="alert">
            <b>{this.state.list.length}</b> Anime fetched!
          </div>
          <ul
            className="list-group">
            {
              this.state.list.map((e, i) => (
                // href={`https://myanimelist.net/anime/${e.id}`}
                // rel="noreferrer"
                <a
                  className="list-group-item list-group-item-action"
                  key={i}
                  onClick={() => this.onAnimeClicked(e)}
                >
                  <span className="name">{e.name}</span>
                  {e.year
                    ? <span className="badge badge-secondary float-right">{e.year}</span>
                    : ''}

                  {e.altNames && e.altNames.length > 0
                    ? <p className="alt">{e.altNames.join(",")}</p>
                    : ''}
                </a>
              ))
            }
          </ul>
        </main>

        <AnimeInfo
          opened={this.state.infoShown}
          anime={this.state.anime}
          onAnimeClosed={this.onAnimeClosed.bind(this)}
        />
      </React.Fragment>
    );
  }

  //#endregion

  //#region Events

  onAnimeClicked(anime) {
    this.setState({ infoShown: true, anime });
  }

  onAnimeClosed() {
    this.setState({ infoShown: false });
  }

  //#endregion
}
